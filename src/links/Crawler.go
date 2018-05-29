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

