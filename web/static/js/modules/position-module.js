import ProductTable from '../components/tables/product-table.js'
import SellerDetails from '../components/seller-details.js'
import ModeSelector from '../components/mode-selector.js'
import PositionTable from '../components/tables/position-table.js'
import APIModule from './api-module.js'

/**
 * Created by RTT.
 * Author: teocci@yandex.com on 2022-8월-29
 */
export default class PositionModule {
    static get instance() {
        this._instance = this._instance ?? new PositionModule()

        return this._instance
    }

    /** @type {SellerDetails} */
    info

    /** @type {ModeSelector} */
    selector

    /** @type {PositionTable} */
    table

    constructor() {
        this.initElement()
        this.initListeners()
    }

    initElement() {
        this.placeholder = document.getElementById('main')

        const $main = this.placeholder

        this.info = new SellerDetails($main)
        this.selector = new ModeSelector($main)
        this.table = new PositionTable($main)
    }

    initListeners() {
        if (this.selector) {
            this.selector.onRequestLoadTable = req => {
                APIModule.requestPositions(req, d => {
                    this.table.updateTable(d)
                })
            }
        }

        if (this.table) {
            this.table.onStateChange = value => {
                switch (value) {
                    case PositionTable.STATE_KEY_DATA_LOADED:
                        this.info.enableButton()
                        break

                    default:
                        this.info.disableButton()
                }
            }
        }

        this.info.$button.onclick = e => {
            this.table.exportTableToXlsx()
        }
    }
}