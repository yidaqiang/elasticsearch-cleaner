package main

import (
	"github.com/spf13/cobra"
	"github.com/yidaqiang/elasticsearch-manager/pkg/action"
	yitypes "github.com/yidaqiang/elasticsearch-manager/pkg/types"
)

const deleteIndicesDesc = ``

func newDeleteIndicesCmd(manager *yitypes.ElasticSearchManager, args []string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete Redundant indices",
		Long:  deleteIndicesDesc,
		Run: func(cmd *cobra.Command, args []string) {
			action.RunDeleteIndices(manager)
		},
	}

	return cmd
}
