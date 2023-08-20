/**
 * Created by RTT.
 * Author: teocci@yandex.com on 2022-12월-22
 */
import APIModule from '../../modules/api-module.js'
import BaseTable from '../../base/base-table.js'

const TAG = 'product-table'

const currencyFormatter = new Intl.NumberFormat('ru-RU', {
    style: 'currency',
    currency: 'RUB',
})

const STATE_KEY_INIT = 'initialized'
const STATE_KEY_DATA_LOADED = 'data-loaded'
const STATE_KEY_DATA_EMPTY = 'data-empty'

const MAIN_SELLER = 25169

const BOX_KEY_SELLER = 'seller'
const BOX_KEY_PRODUCT = 'product'
const BOX_KEY_PRICES = 'prices'

const BOX_SELLER = {
    id: BOX_KEY_SELLER,
    label: 'Seller',
}
const BOX_PRODUCT = {
    id: BOX_KEY_PRODUCT,
    label: 'Product',
}
const BOX_PRICES = {
    id: BOX_KEY_PRICES,
    label: 'Prices',
}

const BOXES = {
    [BOX_KEY_SELLER]: BOX_SELLER,
    [BOX_KEY_PRODUCT]: BOX_PRODUCT,
    [BOX_KEY_PRICES]: BOX_PRICES,
}

const PRICE_KEY_BASE = 'base'
const PRICE_KEY_SALE = 'sale'

const PRICE_BASE = {
    id: PRICE_KEY_BASE,
    label: 'Base',
}
const PRICE_SALE = {
    id: PRICE_KEY_SALE,
    label: 'Sale',
}
const PRICES = {
    [PRICE_KEY_BASE]: PRICE_BASE,
    [PRICE_KEY_SALE]: PRICE_SALE,
}

export default class ProductTable extends BaseTable {
    static TAG = TAG

    constructor(element) {
        super(element)

        this.initProductTableElements()
        this.initProductTableListeners()

        const res = {
            action: 'seller',
            uid: pageInfo.sellerId,
            limit: pageInfo.limit,
        }
        APIModule.fetchProducts(res, d => {
            this.onProductDataFetched(d)
        })
    }

    initProductTableElements() {
        const $wrapper = document.createElement('div')
        $wrapper.classList.add('products', 'module-wrapper')

        this.grid = new gridjs.Grid({
            columns: [
                {
                    name: 'ID',
                    width: '80px',
                },
                {
                    name: 'Name',
                    width: '200px',
                },
                {
                    name: 'Price',
                    width: '220px',
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
                    sort: false,
                    width: '400px',
                    formatter: cell => {
                        console.log({cell})
                        const {identical} = cell

                        if (isNil(identical)) return

                        const $tmp = document.createElement('div')
                        const $wrapper = document.createElement('div')
                        $wrapper.classList.add('identical-wrapper')

                        for (const product of identical) {
                            const boxes = this.createProduct($wrapper)

                            const $seller = boxes.get(BOX_KEY_SELLER).querySelector('.value')
                            this.createBoxInfo($seller, product.supplierInfo, 'sv-item')

                            const $product = boxes.get(BOX_KEY_PRODUCT).querySelector('.value')
                            this.createBoxInfo($product, product, 'pv-item')

                            const $prices = boxes.get(BOX_KEY_PRICES).querySelector('.value')
                            const values = {
                                [PRICE_KEY_BASE]: {
                                    value: product.priceU,
                                    lower: product.priceU < cell.priceU,
                                },
                                [PRICE_KEY_SALE]: {
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
            pagination: {
                limit: 10,
                summary: true,
            },
            data: () => new Promise(resolve => {
                this.resolver = resolve
            }),
        })

        this.grid.render($wrapper)

        this.$head = document.querySelector('.gridjs-head')
        this.$table = document.querySelector('.gridjs-table')

        this.domWithHolderUpdate = $wrapper
    }

    initProductTableListeners() {
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

    onProductDataFetched(d) {
        const data = d.products.map(product => [product.id, product.name, product.priceU / 100, product.salePriceU / 100, product])

        this.data = []
        for (const product of d.products) {
            const {identical} = product
            if (isNil(identical)) continue

            const item = {
                id: product.id,
                name: product.name,
                prices: {
                    base: product.priceU / 100,
                    sale: product.salePriceU / 100,
                },
                competitors: [],
            }

            for (const idem of identical) {
                const competitor = {
                    id: idem.supplierInfo.id,
                    name: idem.supplierInfo.name,
                    base: idem.priceU / 100,
                    sale: idem.salePriceU / 100,
                }
                item.competitors.push(competitor)
            }

            this.data.push(item)
        }

        this.state = data.length > 0 ? STATE_KEY_DATA_LOADED : STATE_KEY_DATA_EMPTY

        this.resolver(data)
    }

    exportTableToXlsx() {
        const sellerId = pageInfo.sellerId
        const data = this.data
        if (isNil(data)) return

        const date = todayToYYYYMMDD()
        const file = `${date}-${sellerId}.xlsx`

        console.log({data})
        const header = ['ID', 'Name', 'Prices', '', 'Competitors', '', '', '']
        const subHeader = ['', '', 'Base', 'Sale', 'ID', 'Name', 'Base', 'Sale']
        const rows = []

        for (const product of data) {
            const row = [
                product.id,
                product.name,
                product.prices.base,
                product.prices.sale,
            ]

            for (const competitor of product.competitors) {
                row.push(
                    competitor.id,
                    competitor.name,
                    competitor.base,
                    competitor.sale,
                )
                rows.push(row.slice()) // create a copy of row to avoid changing the original array
                row.splice(4) // remove the components data from the row
            }
        }

        const worksheet = XLSX.utils.aoa_to_sheet([header, subHeader, ...rows])
        worksheet['!merges'] = [
            {s: {r: 0, c: 0}, e: {r: 1, c: 0}},
            {s: {r: 0, c: 1}, e: {r: 1, c: 1}},
            {s: {r: 0, c: 2}, e: {r: 0, c: 3}},
            {s: {r: 0, c: 4}, e: {r: 0, c: 7}},
        ]

        const workbook = XLSX.utils.book_new()
        XLSX.utils.book_append_sheet(workbook, worksheet, `${sellerId}`)

        return XLSX.writeFile(workbook, file)
    }
}