package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/trilitech/tzgo/rpc"
	"github.com/trilitech/tzgo/tezos"
)

type Params struct {
	BakerAddresses []string
	OutputPath     string
	StartCycle     int64
	EndCycle       int64
	RpcUrl         string
}

func parseParams() *Params {
	rawBakerAddr := flag.String("baker", "", "a list of baker addresses separated by comma")
	outputPath := flag.String("output", "", "output file path")
	rpcUrl := flag.String("rpc", "", "RPC URL")

	startCycle := flag.Int64("start", -1, "start cycle")
	endCycle := flag.Int64("end", -1, "end cycle")

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

	return &Params{
		BakerAddresses: bakers,
		OutputPath:     *outputPath,
		StartCycle:     *startCycle,
		EndCycle:       *endCycle,
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

	head, err := c.GetHeadBlock(ctx)
	if err != nil {
		return
	} else {
		if params.StartCycle == -1 {
			params.StartCycle = head.GetCycle() - 1
		}
		if params.EndCycle == -1 {
			params.EndCycle = head.GetCycle() - 1
		}
		if params.EndCycle < params.StartCycle {
			params.EndCycle = params.StartCycle
		}
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

	fi.WriteString("cycle,start_time,end_time" + strings.Repeat(",total_income,total_loss", len(params.BakerAddresses)) + "\n")
	s := ""
	for cycle := params.StartCycle; cycle <= params.EndCycle; cycle++ {
		startHeight, endHeight := getHeights(head.ChainId, cycle)
		startBlock, err := c.GetBlockHeight(ctx, startHeight)
		if err != nil {
			fmt.Printf("Failed to read start block at level %v: %+v\n", startHeight, err)
			return
		}
		endBlock, err := c.GetBlockHeight(ctx, endHeight)
		if err != nil {
			fmt.Printf("Failed to read end block at level %v: %+v\n", endHeight, err)
			return
		}
		s = fmt.Sprintf("%d,%s,%s", cycle, startBlock.Header.Timestamp.UTC().Format("2006-01-02T15:04:05Z"), endBlock.Header.Timestamp.UTC().Format("2006-01-02T15:04:05Z"))
		for _, baker := range params.BakerAddresses {
			startBalance, err := c.GetDelegate(ctx, tezos.MustParseAddress(baker), startBlock.Hash)
			if err != nil {
				fmt.Printf("Failed to fetch baker %v at level %v: %+v\n", baker, startHeight, err)
				return
			}
			endBalance, err := c.GetDelegate(ctx, tezos.MustParseAddress(baker), endBlock.Hash)
			if err != nil {
				fmt.Printf("Failed to fetch baker %v at level %v: %+v\n", baker, endHeight, err)
				return
			}
			v := float64(endBalance.FullBalance-startBalance.FullBalance) / 1e6
			income := v
			loss := float64(0)
			if v < 0 {
				income = 0
				loss = v
			}
			s += fmt.Sprintf(",%f,%f", income, loss)
		}

		fi.WriteString(s + "\n")
		s = ""
	}
}

func getHeights(chainId tezos.ChainIdHash, cycle int64) (int64, int64) {
	d := tezos.Deployments[chainId].AtCycle(cycle)
	// balance at the end of the last cycle should be the same as that
	// at the very beginning of the current cycle
	startHeight := d.StartHeight + (cycle-d.StartCycle)*d.BlocksPerCycle - 1
	endHeight := startHeight + d.BlocksPerCycle
	return startHeight, endHeight
}
