package main

import (
	monstros "bestiario/internal"
	"context"
	"fmt"

	"github.com/google/uuid"
)

func main() {
	repo := monstros.NewMonstroRepository("root:mulusk.502@tcp(localhost:3306)/bestiario?charset=utf8&parseTime=True&loc=Local")

	repo.Open()

	monstro, err := repo.GetMonstro(context.TODO(), uuid.New())

	if err != nil {
		fmt.Printf("Monstro não encontrado\n")

	} else {
		fmt.Println(monstro)
	}

	novoMonstro := monstros.Monstro{
		UUID:      uuid.New(),
		Nome:      "Lobisomem",
		Descricao: "É bom ter uma arma de prata quando for enfrentar um.",
	}
	monstroSalvo, err := repo.SaveMonstro(context.TODO(), novoMonstro)

	if err != nil {
		panic(err)
	}

	fmt.Printf("Salvo: %v", monstroSalvo)

	monstro, err = repo.SearchMonstros(context.TODO(), "lob")

	if err != nil {
		fmt.Printf("Monstro não encontrado\n")

	} else {
		fmt.Printf("monstro: %v\n", monstro.Nome)
	}
	monstro.Nome = "Lobo"

	repo.UpdateMonstro(context.TODO(), monstro)

	monstro, err = repo.SearchMonstros(context.TODO(), "lob")

	if err != nil {
		fmt.Printf("Monstro não encontrado\n")

	} else {
		fmt.Printf("monstro: %v\n", monstro.Nome)
	}
	resultado := repo.DeleteMonstro(context.TODO(), monstro)

	fmt.Printf("resultado: %v\n", resultado)
}
