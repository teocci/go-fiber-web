import ProductTable from '../components/product-table.js'

/**
 * Created by RTT.
 * Author: teocci@yandex.com on 2022-8월-29
 */
export default class MainModule {
    static get instance() {
        this._instance = this._instance ?? new MainModule()

        return this._instance
    }

    constructor() {
        this.initElement()
        this.initListeners()
    }

    initElement() {
        this.placeholder = document.getElementById('main')
        const $placeholder = this.placeholder

        this.products = new ProductTable($placeholder)
    }

    initListeners() {}

}