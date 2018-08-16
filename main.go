package main

func main() {
	bc := NewBlockChain("")
	defer bc.db.Close()

	cli := Cli{bc}
	cli.Run()
}
