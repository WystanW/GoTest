//匿名函数实现的爬虫
package links

import (
	"fmt"
	"net/http"
	"golang.org/x/net/html"
)

//深度遍历html每个节点
func forEachNode(n *html.Node, pre, post func(n *html.Node)) {
	if pre != nil {
		pre(n)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, pre, post)
	}
	if post != nil {
		post(n)
	}
}

//解析一个url的html页面返回所有链接
func Extract(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("getting %s: %s", url, resp.Status)
	}
	//解析成doc
	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("parsing %s as HTML: %v", url, err)
	}
	//存所有链接
	var links []string
	//visitNode函数
	visitNode := func(n *html.Node) {
		//找标签<a>
		if n.Type == html.ElementNode && n.Data == "a" {
			//遍历标签<a>的所有属性
			for _, a := range n.Attr {
				//找链接
				if a.Key != "href" {
					continue
				}
				//记录链接
				link, err := resp.Request.URL.Parse(a.Val)
				if err != nil {
					continue // ignore bad URLs
				}
				links = append(links, link.String())
			}
		}
	}
	//深度递归,传变量比全局函数好
	forEachNode(doc, visitNode, nil)
	//这个匿名函数
	return links, nil
}

//package main
//
//import "fmt"
//import (
//"links"
//"log"
//"os"
//)
//
//
////f代表爬虫函数,返回一个url里面的所有url,worklist传入要爬的url列表
//func breadthFirst(f func(item string) []string, worklist []string) {
//	//假如现在worklist的yrl都没爬过
//	seen := make(map[string]bool)
//	//只要worklist有url
//	for len(worklist) > 0 {
//		//临时切片保存所有没爬过的url
//		items := worklist
//		//worklist置空,要记录每个url对应页面的所有url
//		worklist = nil
//		//遍历所有的url
//		for _, item := range items {
//			//如果该url没有被爬取,那就记录为这个url已经被爬取了,然后添加到worklist里面代表这些url又要爬取,一直爬
//			if !seen[item] {
//				seen[item] = true
//				//仔细看append源码,接受的是可变参数,也就是可以接收多个string(golang可变参数...)
//				//f(item)返回的是string切片,加个...代表切片里面的所有string
//				worklist = append(worklist, f(item)...)
//			}
//		}
//	}
//}
//
////爬虫函数,爬取一个url,首先打印url然后返回这个url页面里面所有的url
//func crawl(url string) []string {
//	fmt.Println(url)
//	list, err := links.Extract(url)
//	if err != nil {
//		log.Print(err)
//	}
//	return list
//}
//
//func main() {
//	breadthFirst(crawl, os.Args[1:])
//}

