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

        this.grid = new gridjs.Grid({
            columns: ['Name', 'Language', 'Released At', 'Artist'],
            search: true,
            server: {
                url: 'https://catalog.wb.ru/sellers/catalog?appType=1&couponsGeo=12,3,18,15,21&curr=rub&dest=-1257786&emp=0&lang=ru&locale=ru&pricemarginCoeff=1.0&reg=0&regions=80,64,38,4,83,33,68,70,69,30,86,75,40,1,22,66,31,48,110,71&sort=popular&spp=0&supplier=25169',
                then: data => data.products.map(product => [product.id, product.name, product.brand, product.priceU, product.salePriceU]),
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
}