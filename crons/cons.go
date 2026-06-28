package main

import (
	"fmt"

	"gollux/account"

	"github.com/robfig/cron"
)

func settlement() {
	fmt.Println("Running settlements...")
	account.GenerateSettlement()
}

// func agent() {
// 	fmt.Println("Running agent sales summary...")
// 	account.GenerateAgentSales()
// }

func genSales() {

	account.GenerateSales("upline")
	account.GenerateSales("sub_upline")
	account.GenerateSales("direct")
	account.GenerateSales("personal")
}

func main() {
	fmt.Println("Running crons...")
	//settlement()
	//genSales()
	cSales := cron.New()

	cSettle := cron.New()

	err := cSales.AddFunc("0 1 * * *", genSales) // Every day at 1 AM
	if err != nil {
		fmt.Println("Error scheduling task:", err)
		return
	}

	// Schedule task every Saturday at 3 AM
	errb := cSettle.AddFunc("30 1 * * 6", settlement) // Every saturday at 2 AM
	if errb != nil {
		fmt.Println("Error scheduling task:", err)
		return
	}
	cSales.Start()
	fmt.Println("Agent Sales generation. Running every day at 1 AM.")
	cSettle.Start()
	fmt.Println("Settlement generation. Running every Saturday at 1:30 AM.")

	// Keep the program running
	select {}
}
