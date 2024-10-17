package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/trilitech/tzgo/rpc"
	"github.com/trilitech/tzgo/tezos"
)

type Params struct {
	BakerAddresses []string
	OutputPath     string
	StartDate      time.Time
	EndDate        time.Time
	RpcUrl         string
}

func parseParams() *Params {
	rawBakerAddr := flag.String("baker", "", "a list of baker addresses separated by comma")
	outputPath := flag.String("output", "", "output file path")
	rpcUrl := flag.String("rpc", "", "RPC URL")

	startDate := flag.String("start", "", "start date")
	endDate := flag.String("end", "", "end date")

	flag.Parse()

	if len(*rpcUrl) == 0 {
		fmt.Println("No RPC URL given")
		return nil
	}

	bakers := strings.Split(*rawBakerAddr, ",")
	if len(bakers) == 0 {
		fmt.Println("No baker address given")
		return nil
	}

	if len(*outputPath) == 0 {
		fmt.Println("No output file path given")
		return nil
	}

	start := time.Unix(0, 0)
	if *startDate != "" {
		s, err := time.Parse("2006-01-02", *startDate)
		if err != nil {
			fmt.Println("Failed to parse start date")
			return nil
		}
		start = s
	}

	end := time.Now()
	if *endDate != "" {
		e, err := time.Parse("2006-01-02", *endDate)
		if err != nil {
			fmt.Println("Failed to parse end date")
			return nil
		}
		// to cover the full last day
		e = e.Add(time.Hour * 24)
		end = e
	}

	return &Params{
		BakerAddresses: bakers,
		OutputPath:     *outputPath,
		StartDate:      start,
		EndDate:        end,
		RpcUrl:         *rpcUrl,
	}
}

func main() {
	params := parseParams()
	if params == nil {
		return
	}

	// init SDK client
	c, _ := rpc.NewClient(params.RpcUrl, nil)

	// all SDK functions take a context, here we just use a dummy
	ctx := context.TODO()

	startBlock, err := findBlock(c, ctx, params.StartDate)
	if err != nil {
		fmt.Println("Failed to read start block")
		return
	}
	endBlock, err := findBlock(c, ctx, params.EndDate)
	if err != nil {
		fmt.Println("Failed to read end block")
		return
	}

	fi, err := os.Create(params.OutputPath)
	if err != nil {
		fmt.Printf("Failed to create output file: %+v\n", err)
		panic(err)
	}

	defer func() {
		if err := fi.Close(); err != nil {
			fmt.Printf("Failed to close output file: %+v\n", err)
			panic(err)
		}
	}()
	var total float64
	fi.WriteString("baker,start_time,start_block,end_time,end_block,total_income,total_loss\n")
	for _, baker := range params.BakerAddresses {
		startBalance, err := c.GetDelegate(ctx, tezos.MustParseAddress(baker), startBlock.Hash)
		if err != nil {
			fmt.Printf("Failed to fetch baker %v at level %v: %+v\n", baker, startBlock.GetLevel(), err)
			return
		}
		endBalance, err := c.GetDelegate(ctx, tezos.MustParseAddress(baker), endBlock.Hash)
		if err != nil {
			fmt.Printf("Failed to fetch baker %v at level %v: %+v\n", baker, endBlock.GetLevel(), err)
			return
		}
		v := float64(endBalance.FullBalance-startBalance.FullBalance) / 1e6
		income := v
		loss := float64(0)
		if v < 0 {
			income = 0
			loss = v
		}
		s := fmt.Sprintf("%s,%s,%d,%s,%d,%f,%f\n", baker, startBlock.Header.Timestamp.UTC().Format("2006-01-02T15:04:05Z"), startBlock.GetLevel(), endBlock.Header.Timestamp.UTC().Format("2006-01-02T15:04:05Z"), endBlock.GetLevel(), income, loss)
		total += income
		fi.WriteString(s)
	}
	fmt.Printf("total income for %d baker(s): %.6f tez\n", len(params.BakerAddresses), total)
}

func findBlock(c *rpc.Client, ctx context.Context, ts time.Time) (*rpc.Block, error) {
	left := int64(1)
	d, _ := c.GetHeadBlock(ctx)
	right := d.Header.Level

	for {
		if left >= right {
			break
		}
		mid := (left + right) / 2
		block, _ := c.GetBlockHeight(ctx, mid)
		if ts.After(block.GetTimestamp()) {
			left = mid + 1
		} else {
			right = mid
		}
	}
	return c.GetBlockHeight(ctx, left)
}
