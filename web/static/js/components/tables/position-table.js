/**
 * Created by RTT.
 * Author: teocci@yandex.com on 2022-12월-22
 */
import BaseTable from '../../base/base-table.js'

const TAG = 'position-table'

const STATE_KEY_INIT = 'initialized'
const STATE_KEY_DATA_LOADED = 'data-loaded'
const STATE_KEY_DATA_EMPTY = 'data-empty'

const MAIN_SELLER = 25169

const GRID_COLUMNS = [
    {
        name: 'No',
        width: '50px',
    },
    {
        name: 'ID',
        width: '100px',
    },
    {
        name: 'Name',
        width: '200px',
    },
    {
        name: 'Keywords',
        columns: [],
    },
]

const GRID_CONFIG = {
    columns: [],
    search: true,
    sort: true,
    fixedHeader: true,
    width: '1120px',
    height: '500px',
    pagination: {
        limit: 10,
        summary: true,
    },
}

export default class PositionTable extends BaseTable {
    static TAG = TAG

    constructor(element) {
        super(element)

        this.initPositionTableElements()
        this.initPositionTableListeners()
    }

    initPositionTableElements() {
        const $wrapper = document.createElement('div')
        $wrapper.classList.add(TAG, 'module-wrapper')

        this.domWithHolderUpdate = $wrapper
    }

    initPositionTableListeners() {
    }

    updateTable(d) {
        if (isNil(d)) return

        const columns = cloner(GRID_COLUMNS)
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
        const config = cloner(GRID_CONFIG)
        config.columns = columns
        config.data = () => new Promise(resolve => {
            this.resolver = resolve
        })

        if (isNil(this.grid)) {
            this.grid = new gridjs.Grid(config)
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

    addTooltips() {
        const calculateTopPosition = $e => {
            const {left, width, top} = $e.getBoundingClientRect()
            const x = left + width / 2
            const y = top

            return {x, y}
        }

        const calculateBottomPosition = $e => {
            const {left, width, top, height} = $e.getBoundingClientRect()
            const x = width / 2
            const y = height

            return {x, y}
        }

        const showTooltip = (e, config) => {
            let $tooltip = document.querySelector('.tooltip') ?? null
            if (isNil($tooltip)) {
                $tooltip = document.createElement('div')
                document.body.appendChild($tooltip)
            }

            const topMiddle = calculateBottomPosition($tooltip)

            const dX = config.x - topMiddle.x
            const dY = config.y - topMiddle.y

            $tooltip.classList.add('tooltip', 'tooltip-arrow', 'show')
            $tooltip.style.position = 'absolute'
            $tooltip.style.whiteSpace = 'nowrap'
            $tooltip.style.left = `${dX}px`
            $tooltip.style.top = `${dY}px`

            $tooltip.textContent = config.content
        }

        const hideTooltip = () => {
            const $tooltip = document.querySelector('.tooltip')
            if (isNil($tooltip)) return

            $tooltip.classList.remove('show')
        }

        const hList = [...document.querySelectorAll('tr th[data-column-id*="key-"]')]
        for (const $header of hList) {
            if (isNil($header)) continue
            const $content = $header.querySelector('.gridjs-th-content')
            if (isNil($content)) continue

            $header.onmouseover = e => {
                const position = calculateTopPosition($header)
                const config = {
                    x: position.x,
                    y: position.y,
                    content: $content.textContent,
                }
                showTooltip(e, config)
            }
            $header.onmouseleave = e => {
                hideTooltip()
            }
        }
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

    exportTableToXlsx() {
        const data = this.data
        if (isNil(data)) return

        const sellerId = pageInfo.sellerId

        const date = todayToYYYYMMDDHHMM()
        const info = this.computeInfo()

        const file = `${date}-${sellerId}-${info}.xlsx`

        const spaces = []
        for (const k of data.keywords) {
            spaces.push('')
        }

        console.log({data})
        const header = ['No', 'ID', 'Name', 'Keywords', ...spaces]
        const subHeader = ['', '', '', ...data.keywords]
        const rows = []

        for (const product of data.products) {
            const row = [
                product.pos,
                product.id,
                product.name,
                ...product.words.map(k => k.pos),
            ]
            console.log('exportTableToXlsx', {row})
            rows.push(row.slice()) // create a copy of row to avoid changing the original array
        }

        const length = data.keywords.length

        const worksheet = XLSX.utils.aoa_to_sheet([header, subHeader, ...rows])
        worksheet['!merges'] = [
            {s: {r: 0, c: 0}, e: {r: 1, c: 0}},
            {s: {r: 0, c: 1}, e: {r: 1, c: 1}},
            {s: {r: 0, c: 2}, e: {r: 1, c: 2}},
            {s: {r: 0, c: 3}, e: {r: 0, c: length}},
        ]

        const workbook = XLSX.utils.book_new()
        XLSX.utils.book_append_sheet(workbook, worksheet, `${sellerId}`)

        return XLSX.writeFile(workbook, file)
    }
}