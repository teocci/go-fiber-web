/**
 * Created by RTT.
 * Author: teocci@yandex.com on 2023-9월-06
 */
import BaseTable from '../../base/base-table.js'

const TAG = 'datable'

const Grid = gridjs.Grid
const GridHTML = gridjs.html
const GridH = gridjs.h
const RowSelection = gridjs.plugins.selection.RowSelection

const STATUS_UID_INIT = 0
const STATUS_UID_LOADING = 1
const STATUS_UID_LOADED = 2
const STATUS_UID_RENDERED = 3
const STATUS_UID_ERROR = 4

/** @typedef {number | string | boolean | HTMLElement} DTCell */
/** @typedef {Object.<string, DTCell>} DTObject */
/** @typedef {object | string | number | boolean | null | undefined} ComponentChild */

/** @typedef {object} DTColumn
 * @property {?string} [id] - to set the column id
 * @property {function | DTCell | null} [data] - to set the column name
 * @property {string} name - to set the column name
 * @property {?string} [width] - to set the column width
 * @property {?string} [align] - to set the column alignment
 * @property {?string} [sort] - to set the column sorting
 * @property {?boolean} [hidden] - to hide the column
 * @property {?string} [formatter] - to set the column formatter
 * @property {HTMLTableCellElement | function | null} [attributes] - to custom cell attributes
 */

/** @typedef {object} DTStyle
 * @property {?object} [container]
 * @property {?object} [table]
 * @property {?object} [th]
 * @property {?object} [td]
 * @property {?object} [tr]
 * @property {?object} [header]
 * @property {?object} [footer]
 */

/** @typedef {object} DTClassName
 * @property {?string} [container] - to set the container class names
 * @property {?string} [table] - to set the table class names
 * @property {?string} [th] - to set the th class names
 * @property {?string} [td] - to set the td class names
 * @property {?string} [header] - to set the header class names
 * @property {?string} [footer] - to set the footer class names
 * @property {?string} [thead] - to set the thead class names
 * @property {?string} [tbody] - to set the tbody class names
 * @property {?string} [search] - to set the search class names
 * @property {?string} [sort] - to set the sort class names
 * @property {?string} [pagination] - to set the pagination container class names
 * @property {?string} [paginationSummary] - to set the pagination summary class names
 * @property {?string} [paginationButton] - to set the pagination button class names
 * @property {?string} [paginationButtonNext] - to set the pagination button next class names
 * @property {?string} [paginationButtonPrev] - to set the pagination button previous class names
 * @property {?string} [paginationButtonCurrent] - to set the pagination button current class names
 * @property {?string} [loading] - to set the loading container class names
 * @property {?string} [notfound] - to set the empty table container class names
 * @property {?string} [error] - to set the error class names
 */

/** @typedef {object} DTLanguage
 * @property {?object} search
 * @property {?string} search.placeholder
 * @property {?object} sort
 * @property {?string} sort.sortAsc
 * @property {?string} sort.sortDesc
 * @property {?string} pagination
 * @property {?string} pagination.previous
 * @property {?string} pagination.next
 * @property {?function} pagination.navigate
 * @property {?function} pagination.page
 * @property {?string} pagination.showing
 * @property {?string} pagination.of
 * @property {?string} pagination.to
 * @property {?string} pagination.results
 * @property {?string} loading
 * @property {?string} noRecordsFound
 * @property {?string} error
 */

/** @typedef {object} DTServerConfig
 * @property {string} url - to set the server url
 * @property {?string} [method] - to set the request method - default: GET
 * @property {?string} [headers] - to set the request headers
 * @property {?object} [body] - to set the request body
 * @property {?function} [then] - to set the response handler to refine/select attributes
 * @property {?function} [handle] - to handle the response
 * @property {?function} [total] - to set the total records
 */

/** @typedef {object} DTSearchConfig
 * @property {?string} [keyword] - to initiate with a keyword
 * @property {?object} [server] - to enable server integration
 * @property {?number} [debounceTimeout] - to customize searchable fields - default: 1000 (ms)
 * @property {?function} [selector] - to customize searchable fields
 * @property {?function} [ignoreHiddenColumns] - to avoid search inside hidden columns - default: true
 */

/** @typedef {object} DTSortConfig
 * @property {?boolean} [multiColumn] - to enable/disable multi column sort
 * @property {?object} [server] - to enable server integration
 */

/** @typedef {object} DTPaginationConfig
 * @property {?number} [limit] - to set the number of rows per page
 * @property {?number} [page] - to set the initial page
 * @property {?boolean} [summary] - to show/hide the pagination summary
 * @property {?boolean} [nextButton] - to show/hide the next button
 * @property {?boolean} [prevButton] - to show/hide the previous button
 * @property {?boolean} [buttonsCount] - number of buttons to display in the pagination
 * @property {?boolean} [resetPageOnUpdate] - to reset the pagination when table is updated
 * @property {?DTServerConfig} [server] - to enable server integration
 */

/** @typedef {object} DTComponentOptions
 * @property {DTCell[][] | DTObject | Function | null} [data]
 * @property {?HTMLElement} [from]
 * @property {?DTServerConfig} [server] - to enable server integration
 * @property {DTColumn[], string[]} columns
 * @property {?DTStyle} [style]
 * @property {?object} [className]
 * @property {?DTLanguage} [language]
 * @property {?string} width - default: 100%
 * @property {?string} height - default: auto
 * @property {?boolean} selector - default: false
 * @property {?boolean} resolver - default: false
 * @property {?boolean} [autoWidth] - default: true
 * @property {?boolean} [fixedHeader] - default: true
 * @property {boolean | DTSearchConfig | null} [search] - to enable the global search plugin
 * @property {boolean | DTSortConfig | null} [sort] - to enable the sorting plugin
 * @property {boolean | DTPaginationConfig | null} [pagination] - to enable the pagination plugin
 */

export default class DatableComponent extends BaseTable {
    static TAG = TAG

    /** @type {DTComponentOptions} */
    options

    grid

    resolver
    dispatcher

    status = STATUS_UID_INIT

    /**
     * @param {HTMLElement} $element
     * @param {DTComponentOptions} options
     */
    constructor($element, options) {
        super($element)

        this.options = cloner(options)
        console.log({options})

        this.initDatableElements()
    }

    initDatableElements() {
        if (isNil(this.options)) throw new Error('Datable options is required')
        if (isNil(this.options.columns)) throw new Error('Datable columns is required')

        if (this.options.resolver === true) {
            this.options.data = () => new Promise(resolve => {
                this.resolver = resolve
            })
        }

        if (this.options.selector === true) {
            const column = {
                name: GridH('input', {
                    id: 'cb-overall',
                    type: 'checkbox',
                    className: 'gridjs-checkbox',
                    onClick: e => {
                        this.onOverallChange(e.target)
                    },
                }),
                id: 'cb-selector',
                width: '40px',
                sort: false,
                data: row => row.isChecked,
                plugin: {
                    component: RowSelection,
                },
            }
            this.options.columns = [column, ...this.options.columns]
        }

        this.grid = new Grid(this.options)

        this.loadDatable()
    }

    loadDatable() {
        this.grid.render(this.dom)

        this.dispatcher = this.grid.config.store.dispatch

        this.grid.config.store.subscribe(this.onStateChange.bind(this))
    }

    configCBoxes() {
        const list = this.$element.querySelectorAll('input[type="checkbox"].gridjs-checkbox')
        for (const $cb of list) {
            $cb.name = 'dt-cb-group'
            $cb.onchange = e => {
                this.onCBoxChange(list)
            }
        }
    }

    loadData(data) {
        this.data = data
        this.grid.updateConfig({
            data: data,
        })
        this.render()
    }

    render() {
        this.grid.forceRender()
        console.log({grid: this.grid, rowSelection: this.grid.config.store.state.rowSelection})
    }

    updateOverall() {
        if (this.grid.config.store.state?.rowSelection == null) return

        const $overall = document.getElementById('cb-overall')
        if ($overall == null) return

        const length = this.grid.config.data.length

        const list = this.grid.config.store.state.rowSelection?.rowIds
        let checkedCount = list.length ?? 0

        console.log('checkedCount', {checkedCount})

        $overall.checked = checkedCount > 0 && checkedCount === length
        $overall.indeterminate = checkedCount > 0 && checkedCount < length
    }

    checkRow = rowId => state => {
        const rowIds = state.rowSelection?.rowIds || []
        return rowIds.indexOf(rowId) > -1 ? state : {
            ...state,
            rowSelection: {
                rowIds: [rowId, ...rowIds],
            },
        }
    }

    uncheckRow = rowId => state => {
        const rowIds = state.rowSelection?.rowIds || []
        const index = rowIds.indexOf(rowId)

        // rowId doesn't exist
        if (index === -1) return state

        const cloned = [...rowIds]
        cloned.splice(index, 1)

        return {
            ...state,
            rowSelection: {
                rowIds: cloned,
            },
        }
    }

    check(uid, cell) {
        this.dispatcher(this.checkRow(uid))
        cell.isChecked = true
    }

    uncheck(uid, cell) {
        this.dispatcher(this.uncheckRow(uid))
        cell.isChecked = false
    }

    subscribe(fn) {
        this.grid.config.store.subscribe(e => {
            e.status === STATUS_UID_RENDERED && fn(e)
        })
    }

    onOverallChange($overall) {
        if (isNil($overall)) return
        console.log('onOverallChange', {state: this.state})

        const isChecked = $overall.checked
        const data = this.grid.config.data
        const state = this.grid.config.store.state

        for (const index in data) {
            const uid = state.data.rows[index].id
            isChecked ? this.check(uid, data[index]) : this.uncheck(uid, data[index])
        }
    }

    onCBoxChange(e, list) {
        this.updateOverall(list)
    }

    onStateChange(current, prev) {
        if (current.status === prev.status) return
        if (current.status === this.status) return

        this.status = current.status

        if (this.status !== STATUS_UID_RENDERED) return

        this.configCBoxes()
        this.updateOverall()
    }

    startLoading() {
        this.grid.updateConfig({
            data: () => new Promise(resolve => {
                this.resolver = resolve
            })
        })
        this.render()
    }
}