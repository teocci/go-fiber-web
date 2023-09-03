/**
 * Created by RTT.
 * Author: teocci@yandex.com on 2022-12월-22
 */
import BaseTable from '../../base/base-table.js'
import APIModule from '../../modules/api-module.js'

const TAG = 'ads-list-table'

const Grid = gridjs.Grid
const GridHTML = gridjs.html
const RowSelection = gridjs.plugins.selection.RowSelection

const STATE_KEY_INIT = 'initialized'
const STATE_KEY_DATA_LOADED = 'data-loaded'
const STATE_KEY_DATA_EMPTY = 'data-empty'

const STATUS_KEY_READY = 'ready'
const STATUS_KEY_FINISHED = 'finished'
const STATUS_KEY_FAILED = 'failed'
const STATUS_KEY_RUNNING = 'running'
const STATUS_KEY_PAUSED = 'paused'

const STATUS_LABEL_READY = 'Готова к запуску'
const STATUS_LABEL_FINISHED = 'Кампания завершена'
const STATUS_LABEL_FAILED = 'отказался'
const STATUS_LABEL_RUNNING = 'Идут показы'
const STATUS_LABEL_PAUSED = 'Кампания на паузе'

const STATUS_READY = {
    uid: 4,
    key: STATUS_KEY_READY,
    label: STATUS_LABEL_READY,
}

const STATUS_FINISHED = {
    uid: 7,
    key: STATUS_KEY_FINISHED,
    label: STATUS_LABEL_FINISHED,
}

const STATUS_FAILED = {
    uid: 8,
    key: STATUS_KEY_FAILED,
    label: STATUS_LABEL_FAILED,
}

const STATUS_RUNNING = {
    uid: 9,
    key: STATUS_KEY_RUNNING,
    label: STATUS_LABEL_RUNNING,
}

const STATUS_PAUSED = {
    uid: 11,
    key: STATUS_KEY_PAUSED,
    label: STATUS_LABEL_PAUSED,
}

const STATUS_LIST = [
    STATUS_READY,
    STATUS_FINISHED,
    STATUS_FAILED,
    STATUS_RUNNING,
    STATUS_PAUSED,
]

const STATUS_MAP = {
    [STATUS_KEY_READY]: STATUS_READY,
    [STATUS_KEY_FINISHED]: STATUS_FINISHED,
    [STATUS_KEY_FAILED]: STATUS_FAILED,
    [STATUS_KEY_RUNNING]: STATUS_RUNNING,
    [STATUS_KEY_PAUSED]: STATUS_PAUSED,
}

const CAMPAIGN_TYPE_KEY_CATALOG = 'catalog'
const CAMPAIGN_TYPE_KEY_PRODUCT = 'product'
const CAMPAIGN_TYPE_KEY_SEARCH = 'search'
const CAMPAIGN_TYPE_KEY_RECOMMENDATION = 'recommendation'
const CAMPAIGN_TYPE_KEY_AUTO = 'auto'
const CAMPAIGN_TYPE_KEY_SEARCH_CATALOG = 'search-catalog'

const CAMPAIGN_TYPE_LABEL_CATALOG = 'Кампания в каталоге'
const CAMPAIGN_TYPE_LABEL_PRODUCT = 'Кампания в карточке товара'
const CAMPAIGN_TYPE_LABEL_SEARCH = 'Кампания в поиске'
const CAMPAIGN_TYPE_LABEL_RECOMMENDATION = 'Кампания в рекомендациях на главной странице'
const CAMPAIGN_TYPE_LABEL_AUTO = 'Автоматическая кампания'
const CAMPAIGN_TYPE_LABEL_SEARCH_CATALOG = 'Поиск + каталог'

const CAMPAIGN_TYPE_CATALOG = {
    uid: 4,
    key: CAMPAIGN_TYPE_KEY_CATALOG,
    label: CAMPAIGN_TYPE_LABEL_CATALOG,
}

const CAMPAIGN_TYPE_PRODUCT = {
    uid: 5,
    key: CAMPAIGN_TYPE_KEY_PRODUCT,
    label: CAMPAIGN_TYPE_LABEL_PRODUCT,
}

const CAMPAIGN_TYPE_SEARCH = {
    uid: 6,
    key: CAMPAIGN_TYPE_KEY_SEARCH,
    label: CAMPAIGN_TYPE_LABEL_SEARCH,
}

const CAMPAIGN_TYPE_RECOMMENDATION = {
    uid: 7,
    key: CAMPAIGN_TYPE_KEY_RECOMMENDATION,
    label: CAMPAIGN_TYPE_LABEL_RECOMMENDATION,
}

const CAMPAIGN_TYPE_AUTO = {
    uid: 8,
    key: CAMPAIGN_TYPE_KEY_AUTO,
    label: CAMPAIGN_TYPE_LABEL_AUTO,
}

const CAMPAIGN_TYPE_SEARCH_CATALOG = {
    uid: 9,
    key: CAMPAIGN_TYPE_KEY_SEARCH_CATALOG,
    label: CAMPAIGN_TYPE_LABEL_SEARCH_CATALOG,
}

const CAMPAIGN_TYPE_LIST = [
    CAMPAIGN_TYPE_CATALOG,
    CAMPAIGN_TYPE_PRODUCT,
    CAMPAIGN_TYPE_SEARCH,
    CAMPAIGN_TYPE_RECOMMENDATION,
    CAMPAIGN_TYPE_AUTO,
    CAMPAIGN_TYPE_SEARCH_CATALOG,
]

const CAMPAIGN_TYPE_MAP = {
    [CAMPAIGN_TYPE_KEY_CATALOG]: CAMPAIGN_TYPE_CATALOG,
    [CAMPAIGN_TYPE_KEY_PRODUCT]: CAMPAIGN_TYPE_PRODUCT,
    [CAMPAIGN_TYPE_KEY_SEARCH]: CAMPAIGN_TYPE_SEARCH,
    [CAMPAIGN_TYPE_KEY_RECOMMENDATION]: CAMPAIGN_TYPE_RECOMMENDATION,
    [CAMPAIGN_TYPE_KEY_AUTO]: CAMPAIGN_TYPE_AUTO,
    [CAMPAIGN_TYPE_KEY_SEARCH_CATALOG]: CAMPAIGN_TYPE_SEARCH_CATALOG,
}

const buttonLabel = status => {
    switch (status) {
        case 9:
            return 'Stop'
        case 4:
        case 11:
            return 'Start'
        case 7:
            return null
        default:
            throw new Error(`Invalid cell status: ${status}`)
    }
}

const TABLE_COLUMNS = [
    {
        name: '',
        id: 'checkbox',
        sort: false,
        plugin: {
            // install the RowSelection plugins
            component: RowSelection,
        },
        width: '50px',
    },
    {
        name: 'No',
        id: 'no',
        width: '40px',
    },
    {
        name: 'ID',
        id: 'id',
        width: '60px',
    },
    {
        name: 'Name',
        id: 'name',
        width: '150px',
    },
    {
        name: '',
        id: 'type',
        hidden: true,
    },
    {
        name: '',
        id: 'status',
        hidden: true,
    },
    {
        name: 'Type',
        id: 'type_name',
        width: '100px',
    },
    {
        name: 'Status',
        id: 'status_name',
        width: '100px',
    },
    {
        name: 'Actions',
        id: 'actions',
        formatter: cell => {
            const label = buttonLabel(cell.status)
            if (isNil(label)) return

            const $wrapper = document.createElement('div')
            $wrapper.classList.add('identical-wrapper')

            const $btn = document.createElement('button')
            $btn.classList.add('btn', 'btn-primary')
            $btn.dataset.campaignId = cell['campaign_id']
            $btn.textContent = label

            $btn.onclick = e => console.log(`Shows "${cell.status}"`)

            $wrapper.append($btn)

            return GridHTML($wrapper.innerHTML)
        },
    },
]

const statusLabel = uid => {
    for (const item of STATUS_LIST) {
        if (item.uid === uid) return item.label
    }
    return ''
}

const typeLabel = uid => {
    for (const item of CAMPAIGN_TYPE_LIST) {
        if (item.uid === uid) return item.label
    }
    return ''
}

const TABLE_CONFIG = {
    columns: [],
    search: true,
    sort: true,
    fixedHeader: true,
    width: '1120px',
    height: '310px',
}

export default class AdsListTable extends BaseTable {
    static TAG = TAG

    grid

    constructor(element) {
        super(element)

        this.initAdsListTableElements()
        this.initAdsListTableListeners()

        this.requestAdsListData()
    }

    initAdsListTableElements() {
        const $wrapper = document.createElement('div')
        $wrapper.classList.add(TAG, 'module-wrapper')

        const config = cloner(TABLE_CONFIG)
        config.columns = cloner(TABLE_COLUMNS)
        config.data = () => new Promise(resolve => {
            this.resolver = resolve
        })

        this.grid = new gridjs.Grid(config)
        this.grid.render(this.dom)

        this.initHeadAndTable()

        this.domWithHolderUpdate = $wrapper
    }

    initAdsListTableListeners() {
    }

    updateTable(d) {
        if (isNil(d)) return

        const columns = cloner(TABLE_COLUMNS)
        for (const column of columns) {
            if (column.name !== 'Keywords') continue
            let count = 1
            column.columns = d.keywords.map(k => {
                return {
                    id: `key-${count++}`,
                    name: k,
                    width: '200px',
                }
            })
        }
        const config = cloner(TABLE_CONFIG)
        config.columns = columns
        config.data = () => new Promise(resolve => {
            this.resolver = resolve
        })

        if (isNil(this.grid)) {
            this.grid = new Grid(config)
            this.grid.render(this.dom)
        } else {
            this.grid.updateConfig(config)
            this.grid.forceRender()
        }

        console.log('updateTable', {d, resolver: this.resolver})

        setTimeout(() => {
            this.initHeadAndTable()
            this.addTooltips()
            this.loadTableData(d)
            loaderComponent.stopLoader()
        }, 500)
    }

    initHeadAndTable() {
        this.$head = document.querySelector('.gridjs-head')
        this.$table = document.querySelector('.gridjs-table')
    }

    /**
     * @typedef {Object} WordInfo
     * @property {string} word
     * @property {number} avgPos
     * @property {number} count
     * @property {number} wb_count
     * @property {number} total
     * @property {number[]} pos
     *
     *
     * @typedef {Object} ProductItem
     * @property {number} id
     * @property {string} name
     * @property {number} brandId
     * @property {string} brand
     * @property {number} supplierId
     * @property {number} pos
     * @property {WordInfo[]} words
     *
     *
     * @typedef {Object} ProductItems
     * @property {string[]} keywords
     * @property {ProductItem[]} products
     *
     * @param d {ProductItems}
     */
    loadTableData(d) {
        console.log('loadTableData', {d})
        /**
         *
         * @param list {WordInfo[]}
         * @returns {string[]}
         */
        const processPosition = list => {
            const words = []
            for (const k of d.keywords) {
                if (isNil(list)) {
                    words.push('-')
                    continue
                }
                const item = list.find(p => p?.word === k)
                if (isNil(item)) {
                    words.push('-')
                    continue
                }

                words.push(item.avgPos)
            }

            return words
        }

        const data = d.products.map(p => [p.pos, p.id, p.name, ...processPosition(p.words)])
        this.data = {
            keywords: d.keywords,
            products: [],
        }
        for (const product of d.products) {
            if (isNil(product)) continue
            const item = {
                pos: product.pos,
                id: product.id,
                name: product.name,
                words: [],
            }

            for (const k of d.keywords) {
                const info = isNil(product.words) ? null : product.words.find(p => p?.word === k)
                const word = {
                    key: k,
                    pos: isNil(info) || info?.avgPos < 1 ? '-' : info.avgPos,
                }
                item.words.push(word)
            }

            this.data.products.push(item)
        }

        this.state = data.length > 0 ? STATE_KEY_DATA_LOADED : STATE_KEY_DATA_EMPTY

        this.resolver(data)
    }

    computeInfo() {
        const {mode, category, xsubject} = currentData

        console.log('computeInfo', {mode, category, xsubject, currentData})

        const m = isNil(mode) ? '' : mode
        const c = isNil(category) ? '' : isNil(mode) ? category : `-${category}`
        const x = isNil(xsubject) ? isNil(mode) && isNil(category) ? 'все' : '-все' :
            isNil(mode) && isNil(category) ? xsubject : `-${xsubject}`

        return `${m}${c}${x}`
    }

    requestAdsListData() {
        const res = {
            sellerId: pageInfo.sellerId,
        }
        APIModule.requestAdsList(res, d => {
            this.onAdsListDataFetched(d)
        })
    }

    onAdsListDataFetched(d) {
        console.log('onAdsListDataFetched', {d})

        const data = d.map(ad => [
            ad.pos,
            ad.id,
            ad.name,
            ad.type,
            ad.status,
            typeLabel(ad.type),
            statusLabel(ad.status),
            ad,
        ])

        this.resolver(data)
    }
}