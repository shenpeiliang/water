package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"water/util"

	"github.com/fogleman/gg"
	"github.com/spf13/viper"
	"golang.org/x/exp/slog"
)

var (
	Viper = viper.New()
)

// 初始化配置
func init() {
	Viper.SetConfigName("config")
	Viper.SetConfigType("toml")
	Viper.AddConfigPath(".")

	err := Viper.ReadInConfig()
	if err != nil {
		log.Fatalln("读取配置文件错误")
	}
}

// 清空输出目录下的所有文件
func clearOutput(outputDir string) {
	// 创建输出目录
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		os.Mkdir(outputDir, 0777)
	}
	//清空输出目录下的所有文件
	files, err := filepath.Glob(outputDir + "*")
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		os.Remove(file)
	}
}

// 获取输入目录下的所有图片文件
func getFiles(dir string, exts ...string) (files []string, err error) {
	err = filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		//文件的扩展名
		ext := filepath.Ext(path)[1:]
		if util.InSlice(exts, ext) {
			files = append(files, path)
		}
		return nil
	})

	return
}

func main() {
	//目录检查
	outputDir := Viper.GetString("dir.output")
	inputDir := Viper.GetString("dir.input")
	if outputDir == "" || inputDir == "" {
		log.Fatalln("请设置输入输出目录")
	}

	//清空输出目录下的所有文件
	clearOutput(outputDir)

	// 获取输入目录下的所有图片文件
	exts := Viper.GetStringSlice("water.types")
	files, err := getFiles(inputDir, exts...)
	if err != nil {
		log.Fatalln(err)
	}

	//文件数
	count := len(files)

	if count == 0 {
		log.Fatalln("未找到图片文件")
	}

	slog.Info("开始处理图片文件...")

	// 循环处理图片文件
	for index, file := range files {
		fmt.Fprintf(os.Stdout, "\r文件处理中[%d/%d]...", index+1, count)

		err = water(file)
		if err != nil {
			fmt.Println(err)
		}
	}

	fmt.Println()
	slog.Info("处理完成")
}

// 水印处理
func water(path string) (err error) {
	// 打开原始图片
	img, err := gg.LoadImage(path)
	if err != nil {
		return
	}

	//文件名
	fileName := filepath.Base(path)

	// 创建绘图上下文
	dc := gg.NewContextForImage(img)

	// 加载字体
	font, err := gg.LoadFontFace(Viper.GetString("font.path"), Viper.GetFloat64("water.size"))
	if err != nil {
		return
	}
	dc.SetFontFace(font)

	// 设置文字颜色
	dc.SetHexColor(Viper.GetString("water.color"))

	// 获取文字的尺寸
	textWidth, textHeight := dc.MeasureString(Viper.GetString("water.text"))

	// 定义水平和垂直间隔
	horizontalSpacing := Viper.GetFloat64("water.horizontal")
	verticalSpacing := Viper.GetFloat64("water.vertical")

	// 计算需要绘制的文本行和列的数量
	//从0开始，文字的宽度+水平间隔，文字的高度+垂直间隔, 向上取整
	maxX := (float64(dc.Width())) / (float64(textWidth) + horizontalSpacing)
	maxY := (float64(dc.Height()) - verticalSpacing*2) / (float64(textHeight) + verticalSpacing)

	// 循环绘制文本
	for i := 0; i < int(maxX); i++ {
		for j := 0; j < int(maxY); j++ {
			// 计算当前文本的绘制位置，从开始坐标开始
			x := (float64(i) * (float64(textWidth) + horizontalSpacing)) + textWidth/2
			y := (float64(j) * (float64(textHeight) + verticalSpacing)) + textHeight/2 + verticalSpacing

			// 应用旋转（如果需要）
			dc.RotateAbout(gg.Radians(Viper.GetFloat64("water.rotate")), x, y) // 旋转中心为当前文本位置
			dc.DrawString(Viper.GetString("water.text"), x, y)
			dc.RotateAbout(-gg.Radians(Viper.GetFloat64("water.rotate")), x, y) // 旋转回原始状态
		}
	}

	// 保存结果
	ext := filepath.Ext(path)[1:]
	if ext == "png" {
		if err = gg.SavePNG(Viper.GetString("dir.output")+fileName, dc.Image()); err != nil {
			return
		}
	}

	if err = gg.SaveJPG(Viper.GetString("dir.output")+fileName, dc.Image(), 100); err != nil {
		return
	}

	return
}
