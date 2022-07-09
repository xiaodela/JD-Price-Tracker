package saver

import "JD-price-tracker/model"

func ItemListSaver(itemList []model.ItemSku) {
	for _, item := range itemList {
		ItemSaver(item)
	}
}
