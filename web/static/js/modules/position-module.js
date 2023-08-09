import ProductTable from '../components/product-table.js'
import SellerDetail from '../components/seller-detail.js'
import ModeSelector from '../components/mode-selector.js'

/**
 * Created by RTT.
 * Author: teocci@yandex.com on 2022-8월-29
 */
export default class PositionModule {
    static get instance() {
        this._instance = this._instance ?? new PositionModule()

        return this._instance
    }

    /** @type {SellerDetail} */
    seller

    /** @type {ProductTable} */
    products

    constructor() {
        this.initElement()
        this.initListeners()
    }

    initElement() {
        this.placeholder = document.getElementById('main')

        const $main = this.placeholder

        this.seller = new SellerDetail($main)
        this.table = new ModeSelector($main)
        // this.products = new ProductTable($main)
    }

    initListeners() {
        if (!this.products) return

        this.products.onStateChange = value => {
            switch (value) {
                case ProductTable.STATE_DATA_LOADED:
                    this.seller.enableButton()
                    break

                default:
                    this.seller.disableButton()
            }
        }

        this.products.onStateChange = value => {
            switch (value) {
                case ProductTable.STATE_DATA_LOADED:
                    this.seller.enableButton()
                    break

                default:
                    this.seller.disableButton()
            }
        }

        this.seller.$button.onclick = e => {
            this.products.exportTableToXlsx()
        }
    }
}