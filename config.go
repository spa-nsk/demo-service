package main

import (
	"fmt"
	"sync/atomic"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var TimeOutRequest uint64
var TimeOutWork uint64
var CountRequest uint64

func loadConfig() {
	err0 := "Ошибка чтения конфигурационного файла: %w \n"
	err1 := "Ошибка в параметре TimeOutRequest"
	err2 := "Ошибка в параметре TimeOutWork"
	err3 := "Ошибка в параметре CountRequest"

	viper.SetConfigName("config") // имя конфигурационного файла без расширения
	viper.SetConfigType("yaml")   // тип конфигурационного файла (если расширение не указано)
	//viper.AddConfigPath("/etc/demo-service/")   // добавить путь для поиска конфигурационного файла
	//viper.AddConfigPath("$HOME/.demo-service")  //
	viper.AddConfigPath("/opt/demo-service")
	viper.AddConfigPath(".")    // путь для конфигурационного файла текущая папка
	err := viper.ReadInConfig() //
	if err != nil {
		panic(fmt.Errorf(err0, err))
	}

	p1, ok := viper.Get("TimeOutRequest").(int)
	if !ok {
		panic(fmt.Errorf(err1))
	}
	p2, ok := viper.Get("TimeOutWork").(int)
	if !ok {
		panic(fmt.Errorf(err2))
	}
	p3, ok := viper.Get("CountRequest").(int)
	if !ok {
		panic(fmt.Errorf(err3))
	}

	TimeOutRequest = uint64(p1)
	TimeOutWork = uint64(p2)
	CountRequest = uint64(p3)

	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Конфигурационный файл", e.Name, "изменен. Обновление конфигурации")
		err := viper.ReadInConfig() //
		if err != nil {             //
			panic(fmt.Errorf(err0, err))
		}

		p1, ok := viper.Get("TimeOutRequest").(int)
		if !ok {
			panic(fmt.Errorf(err1))
		}
		p2, ok := viper.Get("TimeOutWork").(int)
		if !ok {
			panic(fmt.Errorf(err2))
		}
		p3, ok := viper.Get("CountRequest").(int)
		if !ok {
			panic(fmt.Errorf(err3))
		}
		atomic.StoreUint64(&TimeOutRequest, uint64(p1))
		atomic.StoreUint64(&TimeOutWork, uint64(p2))
		atomic.StoreUint64(&CountRequest, uint64(p3))

	})
	viper.WatchConfig()
}
