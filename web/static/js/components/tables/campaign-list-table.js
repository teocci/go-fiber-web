/**
 * Created by RTT.
 * Author: teocci@yandex.com on 2022-12월-22
 */
import BaseTable from '../../base/base-table.js'
import APIModule from '../../modules/api-module.js'
import DatableComponent from './datable-component.js'

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

const ACTION_KEY_START = 'start'
const ACTION_KEY_STOP = 'stop'
const ACTION_KEY_PAUSE = 'pause'

const ACTION_LABEL_START = 'Запустить'
const ACTION_LABEL_STOP = 'Остановить'
const ACTION_LABEL_PAUSE = 'Приостановить'

const ACTION_START = {
    key: ACTION_KEY_START,
    label: ACTION_LABEL_START,
}

const ACTION_STOP = {
    key: ACTION_KEY_STOP,
    label: ACTION_LABEL_STOP,
}

const ACTION_PAUSE = {
    key: ACTION_KEY_PAUSE,
    label: ACTION_LABEL_PAUSE,
}

const ACTION_LIST = [
    ACTION_START,
    ACTION_STOP,
    ACTION_PAUSE,
]

const ACTION_MAP = {
    [ACTION_KEY_START]: ACTION_START,
    [ACTION_KEY_STOP]: ACTION_STOP,
    [ACTION_KEY_PAUSE]: ACTION_PAUSE,
}

const buttonAction = status => {
    switch (status) {
        case 9:
            return ACTION_PAUSE
        case 4:
        case 11:
            return ACTION_START
        case 7:
            return null
        default:
            throw new Error(`Invalid cell status: ${status}`)
    }
}

const TABLE_COLUMNS = [
    // {
    //     name: '',
    //     id: 'checkbox',
    //     sort: false,
    //     plugin: {
    //         // install the RowSelection plugins
    //         component: RowSelection,
    //     },
    //     width: '50px',
    // },
    {
        name: 'No',
        id: 'no',
        width: '40px',
    },
    {
        name: 'ID',
        id: 'id',
        width: '60px',
        sort: true,
    },
    {
        name: 'Name',
        id: 'name',
        width: '150px',
        sort: true,
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
        width: '80px',
        sort: true,
    },
    {
        name: 'Status',
        id: 'status_name',
        width: '80px',
        sort: true,
    },
    {
        name: 'Actions',
        id: 'actions',
        width: '140px',
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
    fixedHeader: true,
    selector: true,
    resolver: true,
    width: '1120px',
    height: '310px',
}

export default class CampaignListTable extends BaseTable {
    static TAG = TAG

    /** @type {DatableComponent} */
    datable

    constructor(element) {
        super(element)

        this.initCampaignListTableElements()
        this.initCampaignTable()

        this.requestCampaignListData()
    }

    initCampaignListTableElements() {
        const $wrapper = document.createElement('div')
        $wrapper.classList.add(TAG, 'module-wrapper')

        // this.grid = new Grid(config)
        // this.grid.render(this.dom)
        // this.grid.config.store.subscribe(() => {
        //     this.updateTableListeners()
        // })

        this.domWithHolderUpdate = $wrapper
    }

    initCampaignTable() {
        const config = cloner(TABLE_CONFIG)
        config.columns = cloner(TABLE_COLUMNS)

        for (const column of config.columns) {
            if (column.id !== 'actions') continue
            column.formatter = cell => {
                const action = buttonAction(cell.status)
                if (isNil(action)) return

                const $wrapper = document.createElement('div')
                $wrapper.classList.add('identical-wrapper')

                const $btn = document.createElement('button')
                $btn.classList.add('btn', 'btn-primary', 'btn-action')
                $btn.dataset.id = cell['id']
                $btn.dataset.campaignId = cell['campaign_id']
                $btn.dataset.action = action.key
                $btn.textContent = action.label

                $wrapper.append($btn)

                return GridHTML($wrapper.innerHTML)
            }
        }

        this.datable = new DatableComponent(this.dom, config)
        this.datable.subscribe(e => {
            this.updateTableListeners(e)
        })

        this.initHeadAndTable()
    }

    updateTableListeners(e) {
        console.log('updateTableListeners', {e})
        const list = [...document.querySelectorAll('.btn-action')]
        for (const item of list) {
            item.onclick = e => this.onActionClicked(e, item)
        }
    }

    // updateTable(d) {
    //     if (isNil(d)) return
    //
    //     const columns = cloner(TABLE_COLUMNS)
    //     for (const column of columns) {
    //         if (column.name !== 'Keywords') continue
    //         let count = 1
    //         column.columns = d.keywords.map(k => {
    //             return {
    //                 id: `key-${count++}`,
    //                 name: k,
    //                 width: '200px',
    //             }
    //         })
    //     }
    //     const config = cloner(TABLE_CONFIG)
    //     config.columns = columns
    //     config.data = () => new Promise(resolve => {
    //         this.resolver = resolve
    //     })
    //
    //     if (isNil(this.grid)) {
    //         this.grid = new Grid(config)
    //         this.grid.render(this.dom)
    //     } else {
    //         this.grid.updateConfig(config)
    //         this.grid.forceRender()
    //     }
    //
    //     console.log('updateTable', {d, resolver: this.resolver})
    //
    //     setTimeout(() => {
    //         this.initHeadAndTable()
    //         this.addTooltips()
    //         this.loadTableData(d)
    //         loaderComponent.stopLoader()
    //     }, 500)
    // }

    initHeadAndTable() {
        this.$head = document.querySelector('.gridjs-head')
        this.$table = document.querySelector('.gridjs-table')
    }

    requestCampaignListData() {
        const res = {
            sellerId: pageInfo.sellerId,
        }
        APIModule.requestCampaignList(res, d => {
            this.onCampaignListDataFetched(d)
        })
    }

    onCampaignListDataFetched(d) {
        console.log('onCampaignListDataFetched', {d})

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

        this.loadTable(data)
    }

    onActionClicked(e, item) {
        if (isNil(item)) return

        const {id, action} = item.dataset
        console.log('onActionClicked', {id, action})

        this.datable.startLoading()

        const res = {
            sellerId: pageInfo.sellerId,
            campaignId: id,
            action,
        }
        APIModule.requestMarketingControl(res, d => {
            this.onMarketingControlResponse(d)
        })
    }

    loadTable(data) {
        this.datable.loadData(data)
    }

    onMarketingControlResponse(d) {
        console.log('onMarketingControlResponse', {d})
        setTimeout(() => this.requestCampaignListData(), 2000)
    }
}