import ProductTable from '../components/product-table.js'
import SellerDetail from '../components/seller-detail.js'

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
        const $placeholder = this.placeholder

        this.seller = new SellerDetail($placeholder)
        // this.products = new ProductTable($placeholder)
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