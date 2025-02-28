package main

import (
	"archive/zip"
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// 解析命令行参数
	zipPath := flag.String("zip", "", "zip文件的绝对路径")
	flag.Parse()

	if *zipPath == "" {
		fmt.Println("请使用 -zip 参数指定zip文件的绝对路径")
		return
	}

	// 打开源zip文件
	reader, err := zip.OpenReader(*zipPath)
	if err != nil {
		fmt.Printf("打开zip文件失败: %v\n", err)
		return
	}
	defer reader.Close()

	// 找到CodeLocatorPlugin开头的jar文件
	var targetJar *zip.File
	for _, file := range reader.File {
		// 获取文件名（不包含路径）
		fileName := filepath.Base(file.Name)
		if strings.HasPrefix(fileName, "CodeLocatorPlugin") && strings.HasSuffix(fileName, ".jar") {
			targetJar = file
			break
		}
	}

	if targetJar == nil {
		fmt.Println("未找到以CodeLocatorPlugin开头的jar文件")
		return
	}

	// 读取jar文件内容
	jarReader, err := targetJar.Open()
	if err != nil {
		fmt.Printf("打开jar文件失败: %v\n", err)
		return
	}
	defer jarReader.Close()

	// 将jar文件内容读入内存
	jarContent := new(bytes.Buffer)
	if _, err := io.Copy(jarContent, jarReader); err != nil {
		fmt.Printf("读取jar文件内容失败: %v\n", err)
		return
	}

	// 创建新的zip reader来读取jar内容
	jarZipReader, err := zip.NewReader(bytes.NewReader(jarContent.Bytes()), int64(jarContent.Len()))
	if err != nil {
		fmt.Printf("创建jar文件reader失败: %v\n", err)
		return
	}

	// 创建新的jar文件内容
	newJarContent := new(bytes.Buffer)
	jarWriter := zip.NewWriter(newJarContent)

	// 遍历jar中的所有文件
	pluginXmlFound := false
	for _, file := range jarZipReader.File {
		// 获取文件名（不包含路径）
		fileName := filepath.Base(file.Name)
		if fileName == "plugin.xml" {
			// 处理plugin.xml文件
			pluginXmlFound = true
			if err := processPluginXml(file, jarWriter); err != nil {
				fmt.Printf("处理plugin.xml失败: %v\n", err)
				return
			}
		} else {
			// 复制其他文件
			if err := copyFileInZip(file, jarWriter); err != nil {
				fmt.Printf("复制文件失败 %s: %v\n", file.Name, err)
				return
			}
		}
	}

	if !pluginXmlFound {
		fmt.Println("在jar文件中未找到plugin.xml")
		return
	}

	// 关闭jar writer
	if err := jarWriter.Close(); err != nil {
		fmt.Printf("关闭jar writer失败: %v\n", err)
		return
	}

	// 创建新的zip文件
	newZipPath := strings.TrimSuffix(*zipPath, ".zip") + "k2.zip"
	newZipFile, err := os.Create(newZipPath)
	if err != nil {
		fmt.Printf("创建新的zip文件失败: %v\n", err)
		return
	}
	defer newZipFile.Close()

	// 创建新的zip writer
	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	// 将修改后的jar文件写入新的zip（保持原始路径）
	w, err := zipWriter.CreateHeader(&zip.FileHeader{
		Name:     targetJar.Name,  // 保持原始完整路径
		Method:   zip.Deflate,
		Modified: targetJar.Modified,
	})
	if err != nil {
		fmt.Printf("创建jar文件在zip中失败: %v\n", err)
		return
	}
	if _, err := io.Copy(w, bytes.NewReader(newJarContent.Bytes())); err != nil {
		fmt.Printf("写入新的jar文件失败: %v\n", err)
		return
	}

	// 复制其他文件到新的zip
	for _, file := range reader.File {
		if file.Name != targetJar.Name {
			if err := copyFileInZip(file, zipWriter); err != nil {
				fmt.Printf("复制文件失败 %s: %v\n", file.Name, err)
				return
			}
		}
	}

	fmt.Printf("成功创建新的zip文件: %s\n", newZipPath)
}

// 处理plugin.xml文件
func processPluginXml(file *zip.File, zipWriter *zip.Writer) error {
	// 读取原始plugin.xml内容
	rc, err := file.Open()
	if err != nil {
		return fmt.Errorf("打开plugin.xml失败: %v", err)
	}
	defer rc.Close()

	content := new(bytes.Buffer)
	if _, err := io.Copy(content, rc); err != nil {
		return fmt.Errorf("读取plugin.xml内容失败: %v", err)
	}

	// 在文件末尾添加新的扩展配置
	// 在</idea-plugin>之前插入新的扩展
	newContent := strings.Replace(
		content.String(),
		"</idea-plugin>",
		`  <extensions defaultExtensionNs="org.jetbrains.kotlin">
    <supportsKotlinPluginMode supportsK1="true" supportsK2="true" />
  </extensions>
</idea-plugin>`,
		1,
	)

	// 创建新的plugin.xml文件（保持原始路径）
	writer, err := zipWriter.CreateHeader(&zip.FileHeader{
		Name:     file.Name,  // 保持原始完整路径
		Method:   zip.Deflate,
		Modified: file.Modified,
	})
	if err != nil {
		return fmt.Errorf("创建新的plugin.xml失败: %v", err)
	}

	if _, err := writer.Write([]byte(newContent)); err != nil {
		return fmt.Errorf("写入新的plugin.xml内容失败: %v", err)
	}

	return nil
}

// 在zip文件中复制文件
func copyFileInZip(file *zip.File, zipWriter *zip.Writer) error {
	// 创建新的文件头
	writer, err := zipWriter.CreateHeader(&zip.FileHeader{
		Name:     file.Name,  // 保持原始完整路径
		Method:   file.Method,
		Modified: file.Modified,
	})
	if err != nil {
		return err
	}

	// 打开源文件
	reader, err := file.Open()
	if err != nil {
		return err
	}
	defer reader.Close()

	// 复制内容
	_, err = io.Copy(writer, reader)
	return err
}


