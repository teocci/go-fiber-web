import SellerDetails from '../components/seller-details.js'
import AdsListTable from '../components/tables/ads-list-table.js'

/**
 * Created by RTT.
 * Author: teocci@yandex.com on 2022-8월-29
 */
export default class AdsModule {
    static get instance() {
        this._instance = this._instance ?? new AdsModule()

        return this._instance
    }

    /** @type {SellerDetails} */
    info

    /** @type {AdsListTable} */
    table

    constructor() {
        this.initElement()
        this.initListeners()
    }

    initElement() {
        this.placeholder = document.getElementById('main')

        const $main = this.placeholder

        this.info = new SellerDetails($main)
        this.info.hideExportButton()

        this.table = new AdsListTable($main)
    }

    initListeners() {
        if (isNil(this.table)) return

        this.table.onStateChange = value => {
            switch (value) {
                case AdsListTable.STATE_KEY_DATA_LOADED:
                    this.info.enableExportButton()
                    break

                default:
                    this.info.disableExportButton()
            }
        }
    }
}