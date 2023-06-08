package ntomd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"

	"github.com/dstotijn/go-notion"
)

func GetMarkdownStringFromNotionPage(notionKey, pageId *string) (*string, error) {
	if notionKey == nil || pageId == nil {
		return nil, errors.New("'notionKey' and 'pageId' are both required")
	}

	client := notion.NewClient(*notionKey)

	page, err := client.FindPageByID(context.Background(), *pageId)
	if err != nil {
		return nil, fmt.Errorf("(GetMarkdownStringFromNotionPage) client.FindPageById: %v", err)
	}

	body := ""

	switch props := page.Properties.(type) {
	case notion.DatabasePageProperties:
		meta := "---\n"
		for k, v := range props {
			pstr := ParseProperty(v)
			if pstr != nil {
				meta += fmt.Sprintf("%v: %v\n", k, *pstr)
			}
		}
		meta += "---\n\n"
		body += meta
	case notion.PageProperties:
		body += "# " + *ParseRichText(props.Title.Title) + "\n\n"
	}

	children, err := client.FindBlockChildrenByID(context.Background(), *pageId, nil)
	if err != nil {
		return nil, fmt.Errorf("(GetMarkdownStringFromNotionPage) client.FindBlockChildrenByID: %v", err)
	}

	isBuildingList := false
	for _, el := range children.Results {
		isListItemBlock := false
		strToAdd := ""

		switch block := el.(type) {
		case *notion.ParagraphBlock:
			bstr := ParseRichText(block.RichText)
			if bstr != nil {
				strToAdd = *bstr
			}
		case *notion.Heading1Block:
			bstr := ParseRichText(block.RichText)
			if bstr != nil {
				strToAdd = fmt.Sprintf("# %v", *bstr)
			}
		case *notion.Heading2Block:
			bstr := ParseRichText(block.RichText)
			if bstr != nil {
				strToAdd = fmt.Sprintf("## %v", *bstr)
			}
		case *notion.Heading3Block:
			bstr := ParseRichText(block.RichText)
			if bstr != nil {
				strToAdd = fmt.Sprintf("### %v", *bstr)
			}
		case *notion.DividerBlock:
			strToAdd = "---"
		case *notion.BulletedListItemBlock:
			isBuildingList = true
			isListItemBlock = true
			bstr := ParseRichText(block.RichText)
			if bstr != nil {
				strToAdd = fmt.Sprintf("- %v", *bstr)
			}
		case *notion.ImageBlock:
			bstr := ParseImageBlock(block)
			if bstr != nil {
				strToAdd = *bstr
			}
		default:
			log.Println("TRIED TO PARSE UNHANDLED BLOCK!!!", reflect.TypeOf(block))
		}

		if strToAdd != "" {
			// was previously building a list, but it ended so add an extra escape
			if isBuildingList && !isListItemBlock {
				isBuildingList = false
				body += "\n"
			}

			// Its a list item, so only one escape, otherwise add two
			if isListItemBlock {
				body += strToAdd + "\n"
			} else {
				body += strToAdd + "\n\n"
			}
		}
	}

	return &body, nil
}

func ParseProperty(prop notion.DatabasePageProperty) *string {
	if prop.Type == "rich_text" || prop.Type == "name" {
		return ParseRichText(prop.RichText)
	}

	if prop.Type == "url" {
		return prop.URL
	}
	return nil
}

func ParseRichText(richText []notion.RichText) *string {
	str := ""
	for _, el := range richText {
		stri := el.Text.Content
		if el.Annotations != nil {
			if el.Annotations.Italic {
				stri = fmt.Sprintf("_%v_", stri)
			}
			if el.Annotations.Bold {
				stri = fmt.Sprintf("**%v**", stri)
			}
		}
		str += stri
	}
	return &str
}

func ParseImageBlock(block *notion.ImageBlock) *string {
	path := ""
	if block.File != nil {
		path = block.File.URL
	}
	if block.External != nil {
		path = block.External.URL
	}
	if path == "" {
		return nil
	}
	caption := ParseRichText(block.Caption)

	str := fmt.Sprintf("![%v](%v)", *caption, path)
	return &str
}
