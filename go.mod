module github.com/davidforest123/goutil

go 1.19

require (
	cloud.google.com/go/cloudtasks v1.12.3
	cloud.google.com/go/compute v1.23.3
	cloud.google.com/go/storage v1.34.1
	github.com/AdguardTeam/dnsproxy v0.51.0
	github.com/Azure/azure-sdk-for-go v68.0.0+incompatible
	github.com/Azure/azure-sdk-for-go/sdk/storage/azqueue v1.0.0
	github.com/Azure/azure-storage-blob-go v0.15.0
	github.com/Azure/go-autorest/autorest/to v0.4.0
	github.com/ChimeraCoder/anaconda v2.0.0+incompatible
	github.com/Cubox-/libping v0.0.0-20181204104622-3011f76aad09
	github.com/Pallinder/go-randomdata v1.2.0
	github.com/PuerkitoBio/goquery v1.8.0
	github.com/SentimensRG/sigctx v0.0.0-20171003180858-c19b774db63b
	github.com/StevenZack/openurl v0.0.0-20190430065139-b25363f65ff8
	github.com/VividCortex/godaemon v1.0.0
	github.com/adshao/go-binance/v2 v2.3.10
	github.com/aliyun/alibaba-cloud-sdk-go v1.62.584
	github.com/aliyun/aliyun-mns-go-sdk v1.0.2
	github.com/aliyun/aliyun-oss-go-sdk v2.2.6+incompatible
	github.com/anacrolix/utp v0.1.0
	github.com/andreyvit/timerounding v0.8.0
	github.com/antchfx/htmlquery v1.3.0
	github.com/avct/uasurfer v0.0.0-20191028135549-26b5daa857f1
	github.com/aws/aws-sdk-go v1.46.4
	github.com/badoux/checkmail v1.2.1
	github.com/bcampbell/fuzzytime v0.0.0-20191010161914-05ea0010feac
	github.com/beevik/ntp v0.3.0
	github.com/benbjohnson/clock v1.3.0
	github.com/bitly/go-simplejson v0.5.0
	github.com/boombuler/barcode v1.0.1
	github.com/btcsuite/btcd v0.23.4
	github.com/btcsuite/btcd/btcutil v1.1.3
	github.com/cavaliergopher/grab/v3 v3.0.1
	github.com/ccding/go-stun v0.1.4
	github.com/cdipaolo/sentiment v0.0.0-20200617002423-c697f64e7f10
	github.com/chromedp/cdproto v0.0.0-20230716001748-3ed7c525ec8b
	github.com/chromedp/chromedp v0.9.1
	github.com/clarkmcc/go-typescript v0.6.0
	github.com/cretz/bine v0.2.0
	github.com/d5/tengo/v2 v2.13.0
	github.com/domainr/whois v0.1.0
	github.com/dop251/goja v0.0.0-20221118162653-d4bf6fde1b86
	github.com/dxhbiz/codec v0.0.1
	github.com/emersion/go-imap v1.2.1
	github.com/emersion/go-message v0.16.0
	github.com/emirpasic/gods v1.18.1
	github.com/ethereum/go-ethereum v1.10.26
	github.com/extrame/xls v0.0.1
	github.com/fatih/structs v1.1.0
	github.com/frankenbeanies/randhex v0.0.0-20191121050539-48f4de439ea4
	github.com/gabriel-vasile/mimetype v1.4.1
	github.com/getlantern/appdir v0.0.0-20200615192800-a0ef1968f4da
	github.com/getlantern/osversion v0.0.0-20190510010111-432ecec19031
	github.com/getlantern/pac v0.0.0-20161019162755-5534aa917168
	github.com/gin-gonic/gin v1.8.1
	github.com/globalsign/publicsuffix v0.0.0-20220930144140-c2cf46b4ccb7
	github.com/go-errors/errors v1.4.2
	github.com/go-mail/mail v2.3.1+incompatible
	github.com/go-sql-driver/mysql v1.7.0
	github.com/go-telegram-bot-api/telegram-bot-api v4.6.4+incompatible
	github.com/go-vgo/robotgo v0.100.10
	github.com/gocarina/gocsv v0.0.0-20230616125104-99d496ca653d
	github.com/golang/snappy v0.0.4
	github.com/google/go-cmp v0.6.0
	github.com/google/go-github v17.0.0+incompatible
	github.com/google/gopacket v1.1.19
	github.com/google/pprof v0.0.0-20230705174524-200ffdc848b8
	github.com/google/uuid v1.4.0
	github.com/gorilla/websocket v1.5.0
	github.com/goware/urlx v0.3.2
	github.com/h2non/filetype v1.1.3
	github.com/hako/durafmt v0.0.0-20210608085754-5c1018a4e16b
	github.com/jackc/pgx v3.6.2+incompatible
	github.com/jackpal/gateway v1.0.7
	github.com/jbenet/go-base58 v0.0.0-20150317085156-6237cf65f3a6
	github.com/joeguo/tldextract v0.0.0-20220507100122-d83daa6adef8
	github.com/johngb/langreg v0.0.0-20150123211413-5c6abc6d19d2
	github.com/jonhoo/drwmutex v0.0.0-20190519183033-0cffe0733098
	github.com/kardianos/osext v0.0.0-20190222173326-2bc1f35cddc0
	github.com/kbinani/screenshot v0.0.0-20210720154843-7d3a670d8329
	github.com/kenshaw/baseconv v0.1.1
	github.com/klauspost/compress v1.16.6
	github.com/klauspost/pgzip v1.2.5
	github.com/liamcurry/domains v0.0.0-20140814060910-2d799f6e350b
	github.com/libp2p/go-netroute v0.2.1
	github.com/likexian/whois v1.14.4
	github.com/lucazulian/cryptocomparego v0.0.0-20190615070552-deae92f3c4b9
	github.com/mailru/easyjson v0.7.7
	github.com/manifoldco/promptui v0.9.0
	github.com/marcsauter/single v0.0.0-20201009143647-9f8d81240be2
	github.com/markcheno/go-talib v0.0.0-20190307022042-cd53a9264d70
	github.com/mcuadros/go-version v0.0.0-20190830083331-035f6764e8d2
	github.com/melbahja/goph v1.3.0
	github.com/miekg/dns v1.1.55
	github.com/mikaa123/imapmq v0.0.0-20161104140140-bd5a5602fd52
	github.com/mohong122/ip2region v0.0.0-20190505055455-f4ef24f6b03d
	github.com/mssola/user_agent v0.5.3
	github.com/nsf/termbox-go v1.1.1
	github.com/nubunto/tts v0.0.0-20160718193239-d183cb25a053
	github.com/openprovider/ecbrates v0.0.0-20161122034436-f3782097d0a7
	github.com/pariz/gountries v0.1.6
	github.com/pkg/browser v0.0.0-20210911075715-681adbf594b8
	github.com/pkg/errors v0.9.1
	github.com/prestonTao/upnp v0.0.0-20220429011949-f141651daac6
	github.com/quic-go/quic-go v0.36.1
	github.com/r3labs/diff v1.1.0
	github.com/radovskyb/watcher v1.0.7
	github.com/richardlehane/characterize v1.0.0
	github.com/ross-spencer/sfclassic v0.0.0-20190809170605-5e8ee763688f
	github.com/saintfish/chardet v0.0.0-20120816061221-3af4cd4741ca
	github.com/samuel/go-zookeeper v0.0.0-20201211165307-7117e9ea2414
	github.com/satori/go.uuid v1.2.0
	github.com/sekrat/aescrypter v1.0.0
	github.com/sfreiberg/gotwilio v1.0.0
	github.com/shirou/gopsutil/v3 v3.22.10
	github.com/shopspring/decimal v1.3.1
	github.com/sirupsen/logrus v1.8.1
	github.com/songgao/water v0.0.0-20200317203138-2b4b6d7c09d8
	github.com/songtianyi/wechat-go v0.0.0-20220713184122-67e759036893
	github.com/stretchr/testify v1.8.4
	github.com/sttts/color v0.0.0-20141121201513-88cfedb834b6
	github.com/takama/daemon v1.0.0
	github.com/taruti/langdetect v0.0.0-20160316071627-327bfa898307
	github.com/tealeg/xlsx v1.0.5
	github.com/tidwall/gjson v1.14.4
	github.com/tkuchiki/parsetime v0.3.0
	github.com/traefik/yaegi v0.14.3
	github.com/tuotoo/qrcode v0.0.0-20220425170535-52ccc2bebf5d
	github.com/tyler-smith/go-bip32 v1.0.0
	github.com/tyler-smith/go-bip39 v1.1.0
	github.com/ulule/deepcopier v0.0.0-20200430083143-45decc6639b6
	github.com/valyala/fasthttp v1.43.0
	github.com/wcharczuk/go-chart v2.0.1+incompatible
	github.com/willf/bloom v2.0.3+incompatible
	github.com/wumansgy/goEncrypt v1.1.0
	github.com/xtaci/kcp-go v5.4.20+incompatible
	github.com/xtaci/smux v1.5.17
	github.com/yeka/zip v0.0.0-20180914125537-d046722c6feb
	github.com/yl2chen/cidranger v1.0.2
	github.com/zgs225/youdao v1.0.0
	go.mongodb.org/mongo-driver v1.11.0
	golang.org/x/crypto v0.14.0
	golang.org/x/net v0.17.0
	golang.org/x/oauth2 v0.13.0
	golang.org/x/sys v0.13.0
	golang.org/x/text v0.13.0
	gonum.org/v1/gonum v0.12.0
	google.golang.org/api v0.149.0
	google.golang.org/genproto v0.0.0-20231030173426-d783a09b4405
	google.golang.org/grpc v1.59.0
	gopkg.in/headzoo/surf.v1 v1.0.1
	upper.io/db.v3 v3.8.0+incompatible
	xorm.io/xorm v1.3.2
)

require (
	cloud.google.com/go v0.110.9 // indirect
	cloud.google.com/go/accessapproval v1.7.3 // indirect
	cloud.google.com/go/accesscontextmanager v1.8.3 // indirect
	cloud.google.com/go/aiplatform v1.51.2 // indirect
	cloud.google.com/go/analytics v0.21.5 // indirect
	cloud.google.com/go/apigateway v1.6.3 // indirect
	cloud.google.com/go/apigeeconnect v1.6.3 // indirect
	cloud.google.com/go/apigeeregistry v0.8.1 // indirect
	cloud.google.com/go/appengine v1.8.3 // indirect
	cloud.google.com/go/area120 v0.8.3 // indirect
	cloud.google.com/go/artifactregistry v1.14.4 // indirect
	cloud.google.com/go/asset v1.15.2 // indirect
	cloud.google.com/go/assuredworkloads v1.11.3 // indirect
	cloud.google.com/go/automl v1.13.3 // indirect
	cloud.google.com/go/baremetalsolution v1.2.2 // indirect
	cloud.google.com/go/batch v1.6.1 // indirect
	cloud.google.com/go/beyondcorp v1.0.2 // indirect
	cloud.google.com/go/bigquery v1.56.0 // indirect
	cloud.google.com/go/billing v1.17.4 // indirect
	cloud.google.com/go/binaryauthorization v1.7.2 // indirect
	cloud.google.com/go/certificatemanager v1.7.3 // indirect
	cloud.google.com/go/channel v1.17.2 // indirect
	cloud.google.com/go/cloudbuild v1.14.2 // indirect
	cloud.google.com/go/clouddms v1.7.2 // indirect
	cloud.google.com/go/compute/metadata v0.2.3 // indirect
	cloud.google.com/go/contactcenterinsights v1.11.2 // indirect
	cloud.google.com/go/container v1.26.2 // indirect
	cloud.google.com/go/containeranalysis v0.11.2 // indirect
	cloud.google.com/go/datacatalog v1.18.2 // indirect
	cloud.google.com/go/dataflow v0.9.3 // indirect
	cloud.google.com/go/dataform v0.8.3 // indirect
	cloud.google.com/go/datafusion v1.7.3 // indirect
	cloud.google.com/go/datalabeling v0.8.3 // indirect
	cloud.google.com/go/dataplex v1.10.2 // indirect
	cloud.google.com/go/dataproc/v2 v2.2.2 // indirect
	cloud.google.com/go/dataqna v0.8.3 // indirect
	cloud.google.com/go/datastore v1.15.0 // indirect
	cloud.google.com/go/datastream v1.10.2 // indirect
	cloud.google.com/go/deploy v1.14.1 // indirect
	cloud.google.com/go/dialogflow v1.44.2 // indirect
	cloud.google.com/go/dlp v1.10.3 // indirect
	cloud.google.com/go/documentai v1.23.4 // indirect
	cloud.google.com/go/domains v0.9.3 // indirect
	cloud.google.com/go/edgecontainer v1.1.3 // indirect
	cloud.google.com/go/errorreporting v0.3.0 // indirect
	cloud.google.com/go/essentialcontacts v1.6.4 // indirect
	cloud.google.com/go/eventarc v1.13.2 // indirect
	cloud.google.com/go/filestore v1.7.3 // indirect
	cloud.google.com/go/firestore v1.14.0 // indirect
	cloud.google.com/go/functions v1.15.3 // indirect
	cloud.google.com/go/gkebackup v1.3.3 // indirect
	cloud.google.com/go/gkeconnect v0.8.3 // indirect
	cloud.google.com/go/gkehub v0.14.3 // indirect
	cloud.google.com/go/gkemulticloud v1.0.2 // indirect
	cloud.google.com/go/gsuiteaddons v1.6.3 // indirect
	cloud.google.com/go/iam v1.1.4 // indirect
	cloud.google.com/go/iap v1.9.2 // indirect
	cloud.google.com/go/ids v1.4.3 // indirect
	cloud.google.com/go/iot v1.7.3 // indirect
	cloud.google.com/go/kms v1.15.4 // indirect
	cloud.google.com/go/language v1.12.1 // indirect
	cloud.google.com/go/lifesciences v0.9.3 // indirect
	cloud.google.com/go/logging v1.8.1 // indirect
	cloud.google.com/go/longrunning v0.5.3 // indirect
	cloud.google.com/go/managedidentities v1.6.3 // indirect
	cloud.google.com/go/maps v1.5.1 // indirect
	cloud.google.com/go/mediatranslation v0.8.3 // indirect
	cloud.google.com/go/memcache v1.10.3 // indirect
	cloud.google.com/go/metastore v1.13.2 // indirect
	cloud.google.com/go/monitoring v1.16.2 // indirect
	cloud.google.com/go/networkconnectivity v1.14.2 // indirect
	cloud.google.com/go/networkmanagement v1.9.2 // indirect
	cloud.google.com/go/networksecurity v0.9.3 // indirect
	cloud.google.com/go/notebooks v1.11.1 // indirect
	cloud.google.com/go/optimization v1.6.1 // indirect
	cloud.google.com/go/orchestration v1.8.3 // indirect
	cloud.google.com/go/orgpolicy v1.11.3 // indirect
	cloud.google.com/go/osconfig v1.12.3 // indirect
	cloud.google.com/go/oslogin v1.12.1 // indirect
	cloud.google.com/go/phishingprotection v0.8.3 // indirect
	cloud.google.com/go/policytroubleshooter v1.10.1 // indirect
	cloud.google.com/go/privatecatalog v0.9.3 // indirect
	cloud.google.com/go/pubsub v1.33.0 // indirect
	cloud.google.com/go/pubsublite v1.8.1 // indirect
	cloud.google.com/go/recaptchaenterprise/v2 v2.8.2 // indirect
	cloud.google.com/go/recommendationengine v0.8.3 // indirect
	cloud.google.com/go/recommender v1.11.2 // indirect
	cloud.google.com/go/redis v1.13.3 // indirect
	cloud.google.com/go/resourcemanager v1.9.3 // indirect
	cloud.google.com/go/resourcesettings v1.6.3 // indirect
	cloud.google.com/go/retail v1.14.3 // indirect
	cloud.google.com/go/run v1.3.2 // indirect
	cloud.google.com/go/scheduler v1.10.3 // indirect
	cloud.google.com/go/secretmanager v1.11.3 // indirect
	cloud.google.com/go/security v1.15.3 // indirect
	cloud.google.com/go/securitycenter v1.24.1 // indirect
	cloud.google.com/go/servicedirectory v1.11.2 // indirect
	cloud.google.com/go/shell v1.7.3 // indirect
	cloud.google.com/go/spanner v1.51.0 // indirect
	cloud.google.com/go/speech v1.19.2 // indirect
	cloud.google.com/go/storagetransfer v1.10.2 // indirect
	cloud.google.com/go/talent v1.6.4 // indirect
	cloud.google.com/go/texttospeech v1.7.3 // indirect
	cloud.google.com/go/tpu v1.6.3 // indirect
	cloud.google.com/go/trace v1.10.3 // indirect
	cloud.google.com/go/translate v1.9.2 // indirect
	cloud.google.com/go/video v1.20.2 // indirect
	cloud.google.com/go/videointelligence v1.11.3 // indirect
	cloud.google.com/go/vision/v2 v2.7.4 // indirect
	cloud.google.com/go/vmmigration v1.7.3 // indirect
	cloud.google.com/go/vmwareengine v1.0.2 // indirect
	cloud.google.com/go/vpcaccess v1.7.3 // indirect
	cloud.google.com/go/webrisk v1.9.3 // indirect
	cloud.google.com/go/websecurityscanner v1.6.3 // indirect
	cloud.google.com/go/workflows v1.12.2 // indirect
	github.com/AdguardTeam/golibs v0.13.4 // indirect
	github.com/Azure/azure-pipeline-go v0.2.3 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.8.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/azidentity v1.4.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.3.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/billing/armbilling v0.6.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v5 v5.2.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v3 v3.0.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources v1.1.1 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armsubscriptions v1.2.0 // indirect
	github.com/Azure/go-autorest v14.2.0+incompatible // indirect
	github.com/Azure/go-autorest/autorest v0.11.29 // indirect
	github.com/Azure/go-autorest/autorest/adal v0.9.22 // indirect
	github.com/Azure/go-autorest/autorest/date v0.3.0 // indirect
	github.com/Azure/go-autorest/autorest/validation v0.3.1 // indirect
	github.com/Azure/go-autorest/logger v0.2.1 // indirect
	github.com/Azure/go-autorest/tracing v0.6.0 // indirect
	github.com/AzureAD/microsoft-authentication-library-for-go v1.1.1 // indirect
	github.com/ChimeraCoder/tokenbucket v0.0.0-20131201223612-c5a927568de7 // indirect
	github.com/FactomProject/basen v0.0.0-20150613233007-fe3947df716e // indirect
	github.com/FactomProject/btcutilecc v0.0.0-20130527213604-d3a63a5752ec // indirect
	github.com/PuerkitoBio/purell v1.1.1 // indirect
	github.com/PuerkitoBio/urlesc v0.0.0-20170810143723-de5bf2ad4578 // indirect
	github.com/StackExchange/wmi v1.2.1 // indirect
	github.com/StevenZack/tools v1.13.11 // indirect
	github.com/aead/chacha20 v0.0.0-20180709150244-8b13a72661da // indirect
	github.com/aead/poly1305 v0.0.0-20180717145839-3fee0db0b635 // indirect
	github.com/ameshkov/dnscrypt/v2 v2.2.7 // indirect
	github.com/ameshkov/dnsstamps v1.0.3 // indirect
	github.com/anacrolix/envpprof v1.1.1 // indirect
	github.com/anacrolix/missinggo v1.3.0 // indirect
	github.com/anacrolix/missinggo/perf v1.0.0 // indirect
	github.com/anacrolix/missinggo/v2 v2.5.1 // indirect
	github.com/anacrolix/sync v0.4.0 // indirect
	github.com/andybalholm/brotli v1.0.5 // indirect
	github.com/andybalholm/cascadia v1.3.1 // indirect
	github.com/antchfx/xpath v1.2.4 // indirect
	github.com/azr/backoff v0.0.0-20160115115103-53511d3c7330 // indirect
	github.com/blang/semver v3.5.1+incompatible // indirect
	github.com/blend/go-sdk v1.20210918.2 // indirect
	github.com/btcsuite/btcd/btcec/v2 v2.2.0 // indirect
	github.com/btcsuite/btcd/chaincfg/chainhash v1.0.1 // indirect
	github.com/cdipaolo/goml v0.0.0-20210723214924-bf439dd662aa // indirect
	github.com/cheggaaa/pb v2.0.6+incompatible // indirect
	github.com/chromedp/sysutil v1.0.0 // indirect
	github.com/chzyer/readline v1.5.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/daviddengcn/go-colortext v1.0.0 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.0.1 // indirect
	github.com/dlclark/regexp2 v1.7.0 // indirect
	github.com/dustin/go-jsonpointer v0.0.0-20160814072949-ba0abeacc3dc // indirect
	github.com/dustin/gojson v0.0.0-20160307161227-2e71ec9dd5ad // indirect
	github.com/emersion/go-sasl v0.0.0-20200509203442-7bfe0ed36a21 // indirect
	github.com/emersion/go-textwrapper v0.0.0-20200911093747-65d896831594 // indirect
	github.com/extrame/ole2 v0.0.0-20160812065207-d69429661ad7 // indirect
	github.com/facebookgo/stack v0.0.0-20160209184415-751773369052 // indirect
	github.com/fsnotify/fsnotify v1.5.1 // indirect
	github.com/garyburd/go-oauth v0.0.0-20180319155456-bca2e7f09a17 // indirect
	github.com/gen2brain/shm v0.0.0-20200228170931-49f9650110c5 // indirect
	github.com/getlantern/byteexec v0.0.0-20170405023437-4cfb26ec74f4 // indirect
	github.com/getlantern/context v0.0.0-20190109183933-c447772a6520 // indirect
	github.com/getlantern/elevate v0.0.0-20210901195629-ce58359e4d0e // indirect
	github.com/getlantern/errors v1.0.1 // indirect
	github.com/getlantern/filepersist v0.0.0-20160317154340-c5f0cd24e799 // indirect
	github.com/getlantern/golog v0.0.0-20210606115803-bce9f9fe5a5f // indirect
	github.com/getlantern/hex v0.0.0-20190417191902-c6586a6fe0b7 // indirect
	github.com/getlantern/hidden v0.0.0-20190325191715-f02dbb02be55 // indirect
	github.com/getlantern/ops v0.0.0-20190325191751-d70cb0d6f85f // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/go-playground/locales v0.14.0 // indirect
	github.com/go-playground/universal-translator v0.18.0 // indirect
	github.com/go-playground/validator/v10 v10.10.0 // indirect
	github.com/go-sourcemap/sourcemap v2.1.3+incompatible // indirect
	github.com/go-stack/stack v1.8.0 // indirect
	github.com/go-task/slim-sprig v0.0.0-20230315185526-52ccab3ef572 // indirect
	github.com/gobwas/httphead v0.1.0 // indirect
	github.com/gobwas/pool v0.2.1 // indirect
	github.com/gobwas/ws v1.2.1 // indirect
	github.com/goccy/go-json v0.9.11 // indirect
	github.com/gogap/errors v0.0.0-20210818113853-edfbba0ddea9 // indirect
	github.com/gogap/stack v0.0.0-20150131034635-fef68dddd4f8 // indirect
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	github.com/golang-jwt/jwt/v4 v4.5.0 // indirect
	github.com/golang-jwt/jwt/v5 v5.0.0 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/mock v1.6.0 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/google/s2a-go v0.1.7 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.2 // indirect
	github.com/googleapis/gax-go/v2 v2.12.0 // indirect
	github.com/gorilla/schema v1.2.0 // indirect
	github.com/headzoo/surf v1.0.1 // indirect
	github.com/headzoo/ut v0.0.0-20181013193318-a13b5a7a02ca // indirect
	github.com/huandu/xstrings v1.3.1 // indirect
	github.com/ianlancetaylor/demangle v0.0.0-20230524184225-eabc099b10ab // indirect
	github.com/jackc/fake v0.0.0-20150926172116-812a484cc733 // indirect
	github.com/jezek/xgb v1.0.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/cpuid/v2 v2.0.9 // indirect
	github.com/klauspost/reedsolomon v1.9.13 // indirect
	github.com/kr/fs v0.1.0 // indirect
	github.com/kylelemons/godebug v1.1.0 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/lufia/plan9stats v0.0.0-20211012122336-39d0f177ccd0 // indirect
	github.com/lxn/win v0.0.0-20210218163916-a377121e959e // indirect
	github.com/maruel/rs v1.1.0 // indirect
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/mattn/go-ieproxy v0.0.1 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/mattn/go-runewidth v0.0.13 // indirect
	github.com/mattn/go-sqlite3 v1.14.9 // indirect
	github.com/mdp/qrterminal v0.0.0-20180608133721-ba5dc6cf021f // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/montanaflynn/stats v0.7.0 // indirect
	github.com/mxk/go-imap v0.0.0-20150429134902-531c36c3f12d // indirect
	github.com/onsi/ginkgo v1.16.5 // indirect
	github.com/onsi/ginkgo/v2 v2.11.0 // indirect
	github.com/opentracing/opentracing-go v1.2.1-0.20220228012449-10b1cf09e00b // indirect
	github.com/otiai10/gosseract v2.2.1+incompatible // indirect
	github.com/oxtoacart/bpool v0.0.0-20190530202638-03653db5a59c // indirect
	github.com/pelletier/go-toml/v2 v2.0.1 // indirect
	github.com/pkg/sftp v1.13.4 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/power-devops/perfstat v0.0.0-20210106213030-5aafc221ea8c // indirect
	github.com/quic-go/qpack v0.4.0 // indirect
	github.com/quic-go/qtls-go1-19 v0.3.2 // indirect
	github.com/quic-go/qtls-go1-20 v0.2.2 // indirect
	github.com/richardlehane/match v1.0.0 // indirect
	github.com/richardlehane/mscfb v1.0.1 // indirect
	github.com/richardlehane/msoleps v1.0.1 // indirect
	github.com/richardlehane/siegfried v1.7.12-0.20190401190115-99b4106927c7 // indirect
	github.com/richardlehane/xmldetect v1.0.0 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/robotn/gohook v0.31.3 // indirect
	github.com/robotn/xgb v0.0.0-20190912153532-2cb92d044934 // indirect
	github.com/robotn/xgbutil v0.0.0-20190912154524-c861d6f87770 // indirect
	github.com/sekrat/sekrat v1.0.1 // indirect
	github.com/shirou/gopsutil v3.21.10+incompatible // indirect
	github.com/songtianyi/laosj v0.0.0-20180909071508-760f7844583a // indirect
	github.com/songtianyi/rrframework v0.0.0-20180901111106-4caefe307b3f // indirect
	github.com/spaolacci/murmur3 v0.0.0-20180118202830-f09979ecbc72 // indirect
	github.com/syndtr/goleveldb v1.0.1-0.20210819022825-2ae1ddf74ef7 // indirect
	github.com/technoweenie/multipartstreamer v1.0.1 // indirect
	github.com/templexxx/cpufeat v0.0.0-20180724012125-cef66df7f161 // indirect
	github.com/templexxx/xor v0.0.0-20191217153810-f85b25db303b // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.0 // indirect
	github.com/tjfoc/gmsm v1.4.1 // indirect
	github.com/tklauser/go-sysconf v0.3.10 // indirect
	github.com/tklauser/numcpus v0.4.0 // indirect
	github.com/tkuchiki/go-timezone v0.2.2 // indirect
	github.com/ugorji/go/codec v1.2.7 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/vcaesar/gops v0.21.3 // indirect
	github.com/vcaesar/imgo v0.30.0 // indirect
	github.com/vcaesar/keycode v0.10.0 // indirect
	github.com/vcaesar/tt v0.20.0 // indirect
	github.com/willf/bitset v1.1.10 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.1 // indirect
	github.com/xdg-go/stringprep v1.0.3 // indirect
	github.com/xtaci/lossyconn v0.0.0-20200209145036-adba10fffc37 // indirect
	github.com/youmark/pkcs8 v0.0.0-20181117223130-1be2e3e5546d // indirect
	github.com/yusufpapurcu/wmi v1.2.2 // indirect
	github.com/zonedb/zonedb v1.0.3544 // indirect
	go.opencensus.io v0.24.0 // indirect
	golang.org/x/exp v0.0.0-20230713183714-613f0c0eb8a1 // indirect
	golang.org/x/image v0.0.0-20220302094943-723b81ca9867 // indirect
	golang.org/x/mod v0.12.0 // indirect
	golang.org/x/sync v0.4.0 // indirect
	golang.org/x/time v0.3.0 // indirect
	golang.org/x/tools v0.11.0 // indirect
	golang.org/x/xerrors v0.0.0-20220907171357-04be3eba64a2 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20231030173426-d783a09b4405 // indirect
	google.golang.org/genproto/googleapis/bytestream v0.0.0-20231030173426-d783a09b4405 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20231030173426-d783a09b4405 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
	gopkg.in/VividCortex/ewma.v1 v1.1.1 // indirect
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/cheggaaa/pb.v2 v2.0.6 // indirect
	gopkg.in/fatih/color.v1 v1.7.0 // indirect
	gopkg.in/ini.v1 v1.66.2 // indirect
	gopkg.in/mail.v2 v2.3.1 // indirect
	gopkg.in/mattn/go-colorable.v0 v0.0.9 // indirect
	gopkg.in/mattn/go-isatty.v0 v0.0.4 // indirect
	gopkg.in/mattn/go-runewidth.v0 v0.0.3 // indirect
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	rsc.io/qr v0.2.0 // indirect
	xorm.io/builder v0.3.11-0.20220531020008-1bd24a7dc978 // indirect
)
