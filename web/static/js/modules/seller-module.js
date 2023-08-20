import ProductTable from '../components/tables/product-table.js'
import SellerDetails from '../components/seller-details.js'

/**
 * Created by RTT.
 * Author: teocci@yandex.com on 2022-8월-29
 */
export default class SellerModule {
    static get instance() {
        this._instance = this._instance ?? new SellerModule()

        return this._instance
    }

    constructor() {
        this.initElement()
        this.initListeners()
    }

    initElement() {
        this.placeholder = document.getElementById('main')
        const $placeholder = this.placeholder

        this.seller = new SellerDetails($placeholder)
        this.products = new ProductTable($placeholder)
    }

    initListeners() {
        this.products.onStateChange = value => {
            switch (value) {
                case ProductTable.STATE_KEY_DATA_LOADED:
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