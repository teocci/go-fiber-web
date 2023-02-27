/**
 * Created by RTT.
 * Author: teocci@yandex.com on 2022-12월-22
 */
import BaseComponent from '../base/base-component.js'

export default class ProductTable extends BaseComponent {
    static TAG = 'table'

    constructor(element) {
        super(element)

        this.initElements()
        this.initListeners()
    }

    initElements() {
        const $wrapper = document.createElement('div')
        $wrapper.classList.add('products', 'module-wrapper')

        const formatter = new Intl.NumberFormat('ru-RU', {
            style: 'currency',
            currency: 'RUB',
        })

        this.grid = new gridjs.Grid({
            columns: ['ID', 'Name', , {
                name: 'Price',
                columns: [{
                    name: 'Base',
                    formatter: cell => formatter.format(cell),
                }, {
                    name: 'Sale',
                    formatter: cell => formatter.format(cell),
                }],
            }],
            search: true,
            sort: true,
            server: {
                url: '/api/v1/products/seller/25169',
                then: d => d.data.products.map(product => [product.id, product.name, product.priceU / 100, product.salePriceU / 100]),
                handle: res => {
                    // no matching records found
                    if (res.status === 404) return {data: []}
                    if (res.ok) return res.json()

                    throw Error('oh no :(')
                },
            },
        })

        this.grid.render($wrapper)

        this.domWithHolderUpdate = $wrapper
    }

    initListeners() {}

    fetchProducts() {
        fetch('/api/v1/products/seller/25169')
            .then((response) => response.json())
            .then((data) => console.log(data))
    }
}