package handler

import (
	"net/http"

	"strings"

	"github.com/k07me-firefirestyle/engine-v01/article/article"
	miniprop "github.com/k07me-firefirestyle/engine-v01/prop"
	"google.golang.org/appengine"
)

func (obj *ArticleHandler) HandleFind(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	cursor := values.Get("cursor")
	userName := values.Get("userName")
	mode := values.Get("mode")
	tag := []string{}
	props := map[string]string{}
	for k, v := range values {
		if strings.HasPrefix(k, "p-") {
			key := strings.Replace(k, "p-", "", 1)
			props[key] = v[0]
		} else if strings.HasPrefix(k, "t-") {
			tag = append(tag, v[0])
		}
	}
	obj.HandleFindBase(w, r, cursor, userName, props, tag, mode)
}

func (obj *ArticleHandler) HandleFindBase(w http.ResponseWriter, r *http.Request, cursor, //
	userName string, props map[string]string, tags []string, mode string) {
	propObj := miniprop.NewMiniProp()
	ctx := appengine.NewContext(r)
	var foundObj *article.FoundArticles
	//
	//
	//
	manager := obj.GetManager()
	q := manager.NewArtQuery()
	if userName != "" {
		q.WithUserName(ctx, userName)
	}

	if len(props) > 0 {
		q.WithProp(ctx, props)
	}

	if len(tags) > 0 {
		q.WithTags(ctx, tags)
	}

	if mode == "+update" || mode == "update" {
		q.WithUpdatePulas(ctx)
	} else if mode == "-update" {
		q.WithUpdateMinus(ctx)
	}
	q.WithLimitOfFinding(ctx, manager.GetLimitOfFinding())
	foundObj = obj.GetManager().FindArticleFromQuery(ctx, q.GetQuery(), cursor, true)

	propObj.SetPropStringList("", "keys", foundObj.StringIds)
	propObj.SetPropString("", "cursorOne", foundObj.CursorOne)
	propObj.SetPropString("", "cursorNext", foundObj.CursorNext)
	w.Write(propObj.ToJson())
	//}
}
