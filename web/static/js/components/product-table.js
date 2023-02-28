/**
 * Created by RTT.
 * Author: teocci@yandex.com on 2022-12월-22
 */
import BaseComponent from '../base/base-component.js'

const currencyFormatter = new Intl.NumberFormat('ru-RU', {
    style: 'currency',
    currency: 'RUB',
})

const BOX_SUPPLIER = 'supplier'
const BOX_PRODUCT = 'product'
const BOX_PRICES = 'prices'

const BOXES = {
    [BOX_SUPPLIER]: {
        id: BOX_SUPPLIER,
        label: 'Supplier',
    },
    [BOX_PRODUCT]: {
        id: BOX_PRODUCT,
        label: 'Product',
    },
    [BOX_PRICES]: {
        id: BOX_PRICES,
        label: 'Prices',
    },
}

const PRICE_BASE = 'base'
const PRICE_SALE = 'sale'

const PRICES = {
    [PRICE_BASE]: {
        id: PRICE_BASE,
        label: 'Base',
    },
    [PRICE_SALE]: {
        id: PRICE_SALE,
        label: 'Sale',
    },
}

export default class ProductTable extends BaseComponent {
    static TAG = 'table'

    constructor(element) {
        super(element)

        this.initElements()
        this.initListeners()

        this.fetchProducts()
    }

    initElements() {
        const $wrapper = document.createElement('div')
        $wrapper.classList.add('products', 'module-wrapper')

        this.grid = new gridjs.Grid({
            columns: [
                'ID',
                'Name',
                {
                    name: 'Price',
                    columns: [{
                        name: 'Base',
                        formatter: cell => currencyFormatter.format(cell),
                    }, {
                        name: 'Sale',
                        formatter: cell => currencyFormatter.format(cell),
                    }],
                },
                {
                    name: 'Competitors',
                    formatter: cell => {
                        console.log({cell})
                        const {identical} = cell

                        const $tmp = document.createElement('div')
                        const $wrapper = document.createElement('div')
                        $wrapper.classList.add('identical-wrapper')

                        for (const product of identical) {
                            const boxes = this.createProduct($wrapper)

                            const $supplier = boxes.get(BOX_SUPPLIER).querySelector('.value')
                            this.createBoxInfo($supplier, product.supplierInfo, 'sv-item')

                            const $product = boxes.get(BOX_PRODUCT).querySelector('.value')
                            this.createBoxInfo($product, product, 'pv-item')

                            const $prices = boxes.get(BOX_PRICES).querySelector('.value')
                            const values = {
                                [PRICE_BASE]: {
                                    value: product.priceU,
                                    lower: product.priceU < cell.priceU,
                                },
                                [PRICE_SALE]: {
                                    value: product.salePriceU,
                                    lower: product.salePriceU < cell.salePriceU,
                                },
                            }
                            this.createPrices($prices, values)
                        }

                        $tmp.append($wrapper)

                        return gridjs.html($tmp.innerHTML)
                    },
                },
            ],
            search: true,
            sort: true,
            data: [],
            pagination: {
                limit: 50,
                summary: false,
            },
        })

        this.grid.render($wrapper)

        this.domWithHolderUpdate = $wrapper
    }

    initListeners() {}

    fetchProducts() {
        fetch('/api/v1/products/seller/25169')
            .then(res => {
                // no matching records found
                if (res.status === 404) return {data: []}
                if (res.ok) return res.json()

                throw Error('oh no :(')
            })
            .then(d => {
                    const data = d.data.products.map(product => [product.id, product.name, product.priceU / 100, product.salePriceU / 100, product])
                    this.grid.updateConfig({
                        data,
                    })
                    this.grid.forceRender()
                },
            )
    }

    createProduct($wrapper) {
        const boxes = new Map()
        const $element = document.createElement('div')
        $element.classList.add('wb-product')

        for (const id in BOXES) {
            const $box = this.createBox(id)
            $element.append($box)

            boxes.set(id, $box)
        }

        $wrapper.append($element)

        return boxes
    }

    createBox(id) {
        const box = BOXES[id]
        const $box = document.createElement('div')
        $box.classList.add(`${id}-box`, 'wbp-item')

        const $label = document.createElement('div')
        $label.classList.add('label')
        $label.textContent = box.label

        const $value = document.createElement('div')
        $value.classList.add('value')

        $box.append($label, $value)

        return $box
    }

    createBoxInfo($holder, item, style) {
        const $id = document.createElement('div')
        $id.classList.add('id', style)
        $id.textContent = item.id

        const $divider = this.createDivider()
        $divider.classList.add(style)

        const $name = document.createElement('div')
        $name.classList.add('name', style)
        $name.textContent = item.name

        $holder.append($id, $divider, $name)

    }

    createDivider() {
        const $divider = document.createElement('div')
        $divider.classList.add('divider')
        $divider.textContent = '|'

        return $divider
    }

    createPrices($holder, values) {
        const prices = new Map()
        for (const id in PRICES) {
            const $price = this.createPrice(id, values[id])
            $price.classList.add('ppv-item')

            $holder.append($price)

            prices.set(id, $price)
        }

        return prices
    }

    createPrice(id, price) {
        const $price = document.createElement('div')
        $price.classList.add(`p-${id}`)

        const $label = document.createElement('div')
        $label.classList.add('label')
        $label.textContent = PRICES[id].label

        const $value = document.createElement('div')
        $value.classList.add('value')
        if (price.lower) $value.classList.add('lower')
        $value.textContent = currencyFormatter.format(price.value / 100)

        $price.append($label, $value)

        return $price
    }
}