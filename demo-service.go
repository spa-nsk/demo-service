package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"sync/atomic"

	"github.com/PuerkitoBio/goquery"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
)

type ResponseData struct {
	ResponseCount uint64
	TimeResponse  time.Duration
}

const baseYandexURL = "https://yandex.ru/search/touch/?service=www.yandex&ui=webmobileapp.yandex&numdoc=50&lr=213&p=0&text="

type responseStruct struct {
	Error error
	Items []responseItem
}

type responseItem struct {
	Host string
	Url  string
}

// автор парсера parseYandexResponse https://github.com/kkhrychikov/revo-testing/blob/main/serp.go
func parseYandexResponse(response []byte) (res responseStruct) {
	res = responseStruct{Items: make([]responseItem, 0)}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(response))
	if err != nil {
		res.Error = fmt.Errorf("can't create parser for body: %v", err)
		return
	}
	items := doc.Find("div.serp-item")
	items.Each(func(i int, selection *goquery.Selection) {
		_, aExists := selection.Attr("data-fast-name")
		_, cidExists := selection.Attr("data-cid")
		if !selection.Is("div.Label") && !aExists && !selection.Is("span.organic__advLabel") && cidExists {
			link := selection.Find("a.Link").First()

			if link != nil {
				urlStr, _ := link.Attr("href")
				dcStr, _ := link.Attr("data-counter")
				if strings.HasPrefix(urlStr, "https://yandex.ru/turbo/") || strings.Contains(urlStr, "turbopages.org") && dcStr != "" {
					var dc []string
					err := json.Unmarshal([]byte(dcStr), &dc)
					if err != nil || len(dc) < 2 {
						return
					}
					urlStr = dc[1]
				}

				u, err := url.Parse(urlStr)
				if err != nil {
					return
				}

				if u.Host == "" || u.Host == "yabs.yandex.ru" {
					return
				}

				res.Items = append(res.Items, responseItem{
					Host: getRootDomain(u.Host),
					Url:  urlStr,
				})
			}
		}
	})
	return res
}

func getRootDomain(domain string) string {
	domain = strings.ToLower(domain)

	parts := strings.Split(domain, ".")
	if len(parts) < 3 {
		return domain
	}

	if _, ok := tlds[strings.Join(parts[len(parts)-2:], ".")]; ok {
		return strings.Join(parts[len(parts)-3:], ".")
	}

	return strings.Join(parts[len(parts)-2:], ".")
}

var tlds = map[string]struct{}{
	"рф":             {},
	"com.ru":         {},
	"exnet.su":       {},
	"net.ru":         {},
	"org.ru":         {},
	"pp.ru":          {},
	"ru":             {},
	"ru.net":         {},
	"su":             {},
	"aero":           {},
	"asia":           {},
	"biz":            {},
	"com":            {},
	"info":           {},
	"mobi":           {},
	"name":           {},
	"net":            {},
	"org":            {},
	"pro":            {},
	"tel":            {},
	"travel":         {},
	"xxx":            {},
	"adygeya.ru":     {},
	"adygeya.su":     {},
	"arkhangelsk.su": {},
	"balashov.su":    {},
	"bashkiria.ru":   {},
	"bashkiria.su":   {},
	"bir.ru":         {},
	"bryansk.su":     {},
	"cbg.ru":         {},
	"dagestan.ru":    {},
	"dagestan.su":    {},
	"grozny.ru":      {},
	"ivanovo.su":     {},
	"kalmykia.ru":    {},
	"kalmykia.su":    {},
	"kaluga.su":      {},
	"karelia.su":     {},
	"khakassia.su":   {},
	"krasnodar.su":   {},
	"kurgan.su":      {},
	"lenug.su":       {},
	"marine.ru":      {},
	"mordovia.ru":    {},
	"mordovia.su":    {},
	"msk.ru":         {},
	"msk.su":         {},
	"murmansk.su":    {},
	"mytis.ru":       {},
	"nalchik.ru":     {},
	"nalchik.su":     {},
	"nov.ru":         {},
	"nov.su":         {},
	"obninsk.su":     {},
	"penza.su":       {},
	"pokrovsk.su":    {},
	"pyatigorsk.ru":  {},
	"sochi.su":       {},
	"spb.ru":         {},
	"spb.su":         {},
	"togliatti.su":   {},
	"troitsk.su":     {},
	"tula.su":        {},
	"tuva.su":        {},
	"vladikavkaz.ru": {},
	"vladikavkaz.su": {},
	"vladimir.ru":    {},
	"vladimir.su":    {},
	"vologda.su":     {},
	"ad":             {},
	"ae":             {},
	"af":             {},
	"ai":             {},
	"al":             {},
	"am":             {},
	"aq":             {},
	"as":             {},
	"at":             {},
	"aw":             {},
	"ax":             {},
	"az":             {},
	"ba":             {},
	"be":             {},
	"bg":             {},
	"bh":             {},
	"bi":             {},
	"bj":             {},
	"bm":             {},
	"bo":             {},
	"bs":             {},
	"bt":             {},
	"ca":             {},
	"cc":             {},
	"cd":             {},
	"cf":             {},
	"cg":             {},
	"ch":             {},
	"ci":             {},
	"cl":             {},
	"cm":             {},
	"cn":             {},
	"co":             {},
	"co.ao":          {},
	"co.bw":          {},
	"co.ck":          {},
	"co.fk":          {},
	"co.id":          {},
	"co.il":          {},
	"co.in":          {},
	"co.ke":          {},
	"co.ls":          {},
	"co.mz":          {},
	"co.no":          {},
	"co.nz":          {},
	"co.th":          {},
	"co.tz":          {},
	"co.uk":          {},
	"co.uz":          {},
	"co.za":          {},
	"co.zm":          {},
	"co.zw":          {},
	"com.ai":         {},
	"com.ar":         {},
	"com.au":         {},
	"com.bd":         {},
	"com.bn":         {},
	"com.br":         {},
	"com.cn":         {},
	"com.cy":         {},
	"com.eg":         {},
	"com.et":         {},
	"com.fj":         {},
	"com.gh":         {},
	"com.gn":         {},
	"com.gt":         {},
	"com.gu":         {},
	"com.hk":         {},
	"com.jm":         {},
	"com.kh":         {},
	"com.kw":         {},
	"com.lb":         {},
	"com.lr":         {},
	"com.mt":         {},
	"com.mv":         {},
	"com.ng":         {},
	"com.ni":         {},
	"com.np":         {},
	"com.nr":         {},
	"com.om":         {},
	"com.pa":         {},
	"com.pl":         {},
	"com.py":         {},
	"com.qa":         {},
	"com.sa":         {},
	"com.sb":         {},
	"com.sg":         {},
	"com.sv":         {},
	"com.sy":         {},
	"com.tr":         {},
	"com.tw":         {},
	"com.ua":         {},
	"com.uy":         {},
	"com.ve":         {},
	"com.vi":         {},
	"com.vn":         {},
	"com.ye":         {},
	"cr":             {},
	"cu":             {},
	"cx":             {},
	"cz":             {},
	"de":             {},
	"dj":             {},
	"dk":             {},
	"dm":             {},
	"do":             {},
	"dz":             {},
	"ec":             {},
	"ee":             {},
	"es":             {},
	"eu":             {},
	"fi":             {},
	"fo":             {},
	"fr":             {},
	"ga":             {},
	"gd":             {},
	"ge":             {},
	"gf":             {},
	"gg":             {},
	"gi":             {},
	"gl":             {},
	"gm":             {},
	"gp":             {},
	"gr":             {},
	"gs":             {},
	"gy":             {},
	"hk":             {},
	"hm":             {},
	"hn":             {},
	"hr":             {},
	"ht":             {},
	"hu":             {},
	"ie":             {},
	"im":             {},
	"in":             {},
	"in.ua":          {},
	"io ":            {},
	"ir":             {},
	"is":             {},
	"it":             {},
	"je":             {},
	"jo":             {},
	"jp":             {},
	"kg":             {},
	"ki":             {},
	"kiev.ua":        {},
	"kn":             {},
	"kr":             {},
	"ky":             {},
	"kz":             {},
	"li":             {},
	"lk":             {},
	"lt":             {},
	"lu":             {},
	"lv":             {},
	"ly":             {},
	"ma":             {},
	"mc":             {},
	"md":             {},
	"me.uk":          {},
	"mg":             {},
	"mk":             {},
	"mo":             {},
	"mp":             {},
	"ms":             {},
	"mt":             {},
	"mu":             {},
	"mw":             {},
	"mx":             {},
	"my":             {},
	"na":             {},
	"nc":             {},
	"net.cn":         {},
	"nf":             {},
	"ng":             {},
	"nl":             {},
	"no":             {},
	"nu":             {},
	"nz":             {},
	"org.cn":         {},
	"org.uk":         {},
	"pe":             {},
	"ph":             {},
	"pk":             {},
	"pl":             {},
	"pn":             {},
	"pr":             {},
	"ps":             {},
	"pt":             {},
	"re":             {},
	"ro":             {},
	"rs":             {},
	"rw":             {},
	"sd":             {},
	"se":             {},
	"sg":             {},
	"si":             {},
	"sk":             {},
	"sl":             {},
	"sm":             {},
	"sn":             {},
	"so":             {},
	"sr":             {},
	"st":             {},
	"sz":             {},
	"tc":             {},
	"td":             {},
	"tg":             {},
	"tj":             {},
	"tk":             {},
	"tl":             {},
	"tm":             {},
	"tn":             {},
	"to":             {},
	"tt":             {},
	"tw":             {},
	"ua":             {},
	"ug":             {},
	"uk":             {},
	"us":             {},
	"vc":             {},
	"vg":             {},
	"vn":             {},
	"vu":             {},
	"ws":             {},
	"academy":        {},
	"accountant":     {},
	"accountants":    {},
	"actor":          {},
	"adult":          {},
	"africa":         {},
	"agency":         {},
	"airforce":       {},
	"apartments":     {},
	"app":            {},
	"army":           {},
	"art":            {},
	"associates":     {},
	"attorney":       {},
	"auction":        {},
	"audio":          {},
	"auto":           {},
	"band":           {},
	"bank":           {},
	"bar":            {},
	"bargains":       {},
	"bayern":         {},
	"beer":           {},
	"berlin":         {},
	"best":           {},
	"bet":            {},
	"bid":            {},
	"bike":           {},
	"bingo":          {},
	"bio":            {},
	"black":          {},
	"blackfriday":    {},
	"blog":           {},
	"blue":           {},
	"boutique":       {},
	"broker":         {},
	"brussels":       {},
	"build":          {},
	"builders":       {},
	"business":       {},
	"buzz":           {},
	"cab":            {},
	"cafe":           {},
	"cam":            {},
	"camera":         {},
	"camp":           {},
	"capital":        {},
	"car":            {},
	"cards":          {},
	"care":           {},
	"career":         {},
	"careers":        {},
	"cars":           {},
	"casa ":          {},
	"cash":           {},
	"casino":         {},
	"cat":            {},
	"catering":       {},
	"center":         {},
	"ceo":            {},
	"charity":        {},
	"chat":           {},
	"cheap":          {},
	"christmas":      {},
	"church":         {},
	"city":           {},
	"claims":         {},
	"cleaning":       {},
	"click":          {},
	"clinic":         {},
	"clothing":       {},
	"cloud":          {},
	"club":           {},
	"coach":          {},
	"codes":          {},
	"coffee":         {},
	"college":        {},
	"cologne":        {},
	"community":      {},
	"company":        {},
	"computer":       {},
	"condos":         {},
	"construction":   {},
	"consulting":     {},
	"contractors":    {},
	"cooking":        {},
	"cool":           {},
	"coop":           {},
	"country":        {},
	"coupons":        {},
	"courses":        {},
	"credit":         {},
	"creditcard":     {},
	"cricket":        {},
	"cruises":        {},
	"dance":          {},
	"date":           {},
	"dating":         {},
	"deals":          {},
	"degree":         {},
	"delivery":       {},
	"democrat":       {},
	"dental":         {},
	"dentist":        {},
	"desi":           {},
	"design":         {},
	"diamonds":       {},
	"diet":           {},
	"digital":        {},
	"direct":         {},
	"directory":      {},
	"discount":       {},
	"doctor":         {},
	"dog":            {},
	"domains":        {},
	"download":       {},
	"earth":          {},
	"education":      {},
	"email":          {},
	"energy":         {},
	"engineer":       {},
	"engineering":    {},
	"enterprises":    {},
	"equipment":      {},
	"estate":         {},
	"events":         {},
	"exchange":       {},
	"expert":         {},
	"exposed":        {},
	"express":        {},
	"fail":           {},
	"faith":          {},
	"family":         {},
	"fans":           {},
	"farm":           {},
	"fashion":        {},
	"film":           {},
	"finance":        {},
	"financial":      {},
	"fish":           {},
	"fishing":        {},
	"fit":            {},
	"fitness":        {},
	"flights":        {},
	"florist":        {},
	"flowers":        {},
	"fm":             {},
	"football":       {},
	"forex":          {},
	"forsale":        {},
	"foundation":     {},
	"fun":            {},
	"fund":           {},
	"furniture":      {},
	"futbol":         {},
	"fyi":            {},
	"gallery":        {},
	"game":           {},
	"games":          {},
	"garden":         {},
	"gent":           {},
	"gift":           {},
	"gifts":          {},
	"gives":          {},
	"glass":          {},
	"global":         {},
	"gmbh":           {},
	"gold":           {},
	"golf":           {},
	"graphics":       {},
	"gratis":         {},
	"green":          {},
	"gripe":          {},
	"group":          {},
	"guide":          {},
	"guitars":        {},
	"guru":           {},
	"haus":           {},
	"healthcare":     {},
	"help":           {},
	"hiphop":         {},
	"hockey":         {},
	"holdings":       {},
	"holiday":        {},
	"horse":          {},
	"hospital":       {},
	"host":           {},
	"hosting":        {},
	"house":          {},
	"how":            {},
	"immo":           {},
	"immobilien":     {},
	"industries":     {},
	"ink":            {},
	"institute":      {},
	"insure":         {},
	"international":  {},
	"investments":    {},
	"irish":          {},
	"jewelry":        {},
	"jobs":           {},
	"juegos":         {},
	"kaufen":         {},
	"kim":            {},
	"kitchen":        {},
	"kiwi":           {},
	"land":           {},
	"lawyer":         {},
	"lease":          {},
	"legal":          {},
	"life":           {},
	"lighting":       {},
	"limited":        {},
	"limo":           {},
	"link":           {},
	"live":           {},
	"llc":            {},
	"loan":           {},
	"loans":          {},
	"lol":            {},
	"london":         {},
	"love":           {},
	"ltd":            {},
	"luxe":           {},
	"luxury":         {},
	"maison":         {},
	"management":     {},
	"market":         {},
	"marketing":      {},
	"mba":            {},
	"media":          {},
	"memorial":       {},
	"men":            {},
	"menu":           {},
	"miami":          {},
	"moda":           {},
	"moe":            {},
	"mom":            {},
	"money":          {},
	"mortgage":       {},
	"moscow":         {},
	"movie":          {},
	"navy":           {},
	"network":        {},
	"news":           {},
	"ninja":          {},
	"observer":       {},
	"one":            {},
	"onl":            {},
	"online":         {},
	"ooo":            {},
	"page":           {},
	"paris":          {},
	"partners":       {},
	"parts":          {},
	"party":          {},
	"pet":            {},
	"photo":          {},
	"photography":    {},
	"photos":         {},
	"pics":           {},
	"pictures":       {},
	"pink":           {},
	"pizza":          {},
	"plumbing":       {},
	"plus":           {},
	"poker":          {},
	"press":          {},
	"productions":    {},
	"promo":          {},
	"properties":     {},
	"property":       {},
	"protection":     {},
	"pub":            {},
	"qpon":           {},
	"racing":         {},
	"radio":          {},
	"radio.am":       {},
	"radio.fm":       {},
	"realty":         {},
	"recipes":        {},
	"red":            {},
	"rehab":          {},
	"reisen":         {},
	"rent":           {},
	"rentals":        {},
	"repair":         {},
	"report":         {},
	"republican":     {},
	"rest":           {},
	"restaurant":     {},
	"review":         {},
	"reviews":        {},
	"rich":           {},
	"rip":            {},
	"rocks":          {},
	"rodeo":          {},
	"run":            {},
	"sale":           {},
	"salon":          {},
	"sarl":           {},
	"school":         {},
	"schule":         {},
	"science":        {},
	"security":       {},
	"services":       {},
	"sex":            {},
	"sexy":           {},
	"shiksha":        {},
	"shoes":          {},
	"shop":           {},
	"shopping":       {},
	"show":           {},
	"singles":        {},
	"site":           {},
	"ski":            {},
	"soccer":         {},
	"social":         {},
	"software":       {},
	"solar":          {},
	"solutions":      {},
	"soy":            {},
	"space":          {},
	"sport":          {},
	"store":          {},
	"stream":         {},
	"studio":         {},
	"study":          {},
	"style":          {},
	"sucks":          {},
	"supplies":       {},
	"supply":         {},
	"support":        {},
	"surf":           {},
	"surgery":        {},
	"systems":        {},
	"tatar":          {},
	"tattoo":         {},
	"tax":            {},
	"taxi":           {},
	"team":           {},
	"tech":           {},
	"technology":     {},
	"tennis":         {},
	"theater":        {},
	"theatre":        {},
	"tickets":        {},
	"tienda":         {},
	"tips":           {},
	"tires":          {},
	"tirol":          {},
	"today":          {},
	"tools":          {},
	"top":            {},
	"tours":          {},
	"town":           {},
	"toys":           {},
	"trade":          {},
	"trading":        {},
	"training":       {},
	"tube":           {},
	"tv":             {},
	"university":     {},
	"uno":            {},
	"vacations":      {},
	"vegas":          {},
	"ventures":       {},
	"vet":            {},
	"viajes":         {},
	"video":          {},
	"villas":         {},
	"vin":            {},
	"vip":            {},
	"vision":         {},
	"vodka":          {},
	"vote":           {},
	"voting":         {},
	"voto":           {},
	"voyage":         {},
	"watch":          {},
	"webcam":         {},
	"website":        {},
	"wedding":        {},
	"wien":           {},
	"wiki":           {},
	"win":            {},
	"wine":           {},
	"work":           {},
	"works":          {},
	"world":          {},
	"wtf":            {},
	"xyz":            {},
	"yoga":           {},
	"zone":           {},
	"дети":           {},
	"москва":         {},
	"онлайн":         {},
	"орг":            {},
	"рус":            {},
	"сайт":           {},
}

func checkAvailability(url string) (uint64, time.Duration) {
	var i, index uint64
	countRequest := atomic.LoadUint64(&CountRequest)
	timeOutRequest := time.Millisecond * time.Duration(atomic.LoadUint64(&TimeOutRequest))
	timeResponse := time.Millisecond * 0

	ch := make(chan time.Duration)

	for i = 0; i < countRequest; i++ {
		go readUrl(url, timeOutRequest, ch)
	}

	for i = 0; i < countRequest; i++ {
		t := <-ch
		if t == time.Second*9999 {
			if index == 0 {
				index = i
			}
			continue
		}
		if t > timeResponse {
			timeResponse = t
		}
	}
	if index == 0 {
		return i, timeResponse
	}
	return index, timeResponse
}

func readUrl(url string, sec time.Duration, ch chan time.Duration) {
	var defaultTtransport http.RoundTripper = &http.Transport{Proxy: nil}
	client := &http.Client{Timeout: sec, Transport: defaultTtransport}
	start := time.Now()
	resp, err := client.Get(url)

	if err != nil {
		ch <- time.Second * 9999
		return
	}
	if resp.StatusCode != 200 {
		ch <- time.Second * 9999
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	_ = body
	if err != nil {
		ch <- time.Second * 9999
		return
	}
	end := time.Now()
	ch <- end.Sub(start)
}

var TimeOutRequest uint64
var TimeOutWork uint64
var CountRequest uint64

func main() {

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

	mux := http.NewServeMux()
	mux.HandleFunc("/sites", searchSites)
	mux.HandleFunc("/sitesclient", clientSearchSites)

	log.Println("Слушаем порт :8080...")
	http.ListenAndServe(":8080", mux)
}

func clientSearchSites(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, http.StatusText(405), 405)
		return
	}
	search := r.URL.Query().Get("search")
	if search == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}
	var defaultTtransport http.RoundTripper = &http.Transport{Proxy: nil}
	client := &http.Client{Transport: defaultTtransport}

	resp, err := client.Get("http://127.0.0.1:8080/sites?search=" + search)

	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	//декодировать ответ и сформировать страничку ответа в структурированном виде
	var s map[string]ResponseData
	err = json.Unmarshal(body, &s)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	tmpl, _ := template.ParseFiles("view/search.html")
	tmpl.Execute(w, &s)
}

func searchSites(w http.ResponseWriter, r *http.Request) {
	timeOutRequest := time.Millisecond * time.Duration(atomic.LoadUint64(&TimeOutWork))
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), timeOutRequest)
	defer cancel()
	defer func() {
		end := time.Now()
		fmt.Println("Время выполнения запроса", end.Sub(start))
	}()
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, http.StatusText(405), 405)
		return
	}

	search := r.URL.Query().Get("search")
	if search == "" {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	var defaultTtransport http.RoundTripper = &http.Transport{Proxy: nil}
	client := &http.Client{Transport: defaultTtransport}

	resp, err := client.Get(baseYandexURL + search)

	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	//func parseYandexResponse(response []byte) (res responseStruct)
	res := parseYandexResponse(body)

	if res.Error != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	s := make(map[string]ResponseData)

	for _, item := range res.Items {
		select {
		case <-ctx.Done():
			fmt.Println("Истекло время выполнения запроса (", timeOutRequest, ").")
			json.NewEncoder(w).Encode(s)
			return
		default:
			count, timeResponse := checkAvailability(item.Url)
			s[item.Host] = ResponseData{count, timeResponse}
			fmt.Println(item.Host, count, timeResponse)
		}
	}
	json.NewEncoder(w).Encode(s)
}
