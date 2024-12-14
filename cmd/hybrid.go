package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/xackery/quail/os"
	"github.com/xackery/quail/quail"
)

func init() {
	rootCmd.AddCommand(hybridCmd)
}

// hybridCmd represents the hybrid command
var hybridCmd = &cobra.Command{
	Use:   "hybrid",
	Short: "Hybrid merge geometry with existing bone data",
	Long:  `Hybrid is a temporary hack while animations and bone data isn't yet implemented.`,
	Run:   runHybrid,
}

func runHybrid(cmd *cobra.Command, args []string) {
	err := runHybridE(cmd, args)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func runHybridE(cmd *cobra.Command, args []string) error {
	if len(args) < 3 {
		fmt.Println("Expected 4 arguments, got", len(args))
		fmt.Println("Usage: quail hybrid <geometry.eqg> <bone.eqg> <out.eqg>")
		os.Exit(1)
	}
	srcPath, err := filepath.Abs(args[0])
	if err != nil {
		return fmt.Errorf("parse source path: %w", err)
	}

	bonePath, err := filepath.Abs(args[1])
	if err != nil {
		return fmt.Errorf("parse bone path: %w", err)
	}

	outPath, err := filepath.Abs(args[2])
	if err != nil {
		return fmt.Errorf("parse out path: %w", err)
	}

	err = hybrid(srcPath, bonePath, outPath)
	if err != nil {
		return fmt.Errorf("hybrid: %w", err)
	}

	return nil
}

func hybrid(srcPath, bonePath, outPath string) error {
	var err error
	if filepath.Ext(srcPath) != ".eqg" {
		return fmt.Errorf("source path must be an eqg")
	}
	if filepath.Ext(bonePath) != ".eqg" {
		return fmt.Errorf("bone path must be an eqg")
	}
	if filepath.Ext(outPath) != ".eqg" {
		return fmt.Errorf("out path must be an eqg")
	}

	_, err = os.Stat(srcPath)
	if err != nil {
		return fmt.Errorf("stat %s: %w", srcPath, err)
	}

	_, err = os.Stat(bonePath)
	if err != nil {
		return fmt.Errorf("stat %s: %w", bonePath, err)
	}

	srcQ := quail.New()
	err = srcQ.PfsRead(srcPath)
	if err != nil {
		return fmt.Errorf("src pfs read %s: %w", srcPath, err)
	}

	boneQ := quail.New()
	err = boneQ.PfsRead(bonePath)
	if err != nil {
		return fmt.Errorf("bone pfs read %s: %w", bonePath, err)
	}

	boneQ.Models = srcQ.Models
	err = boneQ.PfsWrite(1, 1, outPath)
	if err != nil {
		return fmt.Errorf("bone pfs export %s: %w", outPath, err)
	}

	return nil
}
