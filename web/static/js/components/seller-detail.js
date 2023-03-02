/**
 * Created by Teocci.
 * Author: teocci@yandex.com on 2023-Mar-01
 */
import BaseComponent from '../base/base-component.js'

export default class SellerDetail extends BaseComponent {
    static TAG = 'seller-detail'

    constructor(element) {
        super(element)

        this.seller = null

        this.initElements()
        this.initListeners()

        this.fetchSeller()
    }

    initElements() {
        const $component = document.createElement('div')
        $component.classList.add('seller-details', 'component-wrapper')

        const $infoWrap = document.createElement('div')
        $infoWrap.classList.add('info-wrap', 'cw-item')

        const $logoWrap = document.createElement('div')
        $logoWrap.classList.add('logo-wrap')

        const $logo = document.createElement('div')
        $logo.classList.add('logo')

        const $img = document.createElement('img')
        $img.src = '/img/seller-empty-logo.jpg'
        $img.width = '120'
        $img.height = '50'

        const $info = document.createElement('div')
        $info.classList.add('info')

        const $nameWrap = document.createElement('div')
        $nameWrap.classList.add('info-name', 'item-wrap')

        const $name = document.createElement('h2')
        $name.classList.add('name')

        const $trademarkWrap = document.createElement('div')
        $trademarkWrap.classList.add('info-trademark', 'item-wrap')

        const $trademark = document.createElement('div')
        $trademark.classList.add('trademark')

        const $legalAddressWrap = document.createElement('div')
        $legalAddressWrap.classList.add('info-legal-address', 'item-wrap')

        const $legalAddress = document.createElement('div')
        $legalAddress.classList.add('legal-address')

        const $btnWrap = document.createElement('div')
        $btnWrap.classList.add('btn-wrap', 'cw-item')

        const $btn = document.createElement('button')
        $btn.classList.add('export-btn')
        $btn.textContent = 'Export'
        $btn.disabled = true

        $logo.append($img)
        $logoWrap.append($logo)

        $nameWrap.append($name)
        $trademarkWrap.append($trademark)
        $legalAddressWrap.append($legalAddress)

        $info.append($nameWrap, $trademarkWrap, $legalAddressWrap)

        $infoWrap.append($logoWrap, $info)
        $btnWrap.append($btn)
        $component.append($infoWrap, $btnWrap)

        this.$logo = $img
        this.$name = $name
        this.$trademark = $trademark
        this.$legalAddress = $legalAddress
        this.$button = $btn

        this.domWithHolderUpdate = $component
    }

    initListeners() {}

    fetchSeller() {
        const supplierId = pageInfo.supplierId
        const url = `/api/v1/seller/${supplierId}`
        fetch(url)
            .then(res => res.json())
            .then(d => {
                    const seller = d.data
                    this.loadData(seller)
                },
            )
    }

    loadData(seller) {
        this.seller = seller

        this.$name.textContent = seller.name
        this.$trademark.textContent = seller.trademark
        this.$legalAddress.textContent = seller.legalAddress

        this.$logo.src = `/seller/logo/${seller.id}`
        this.$logo.alt = seller.fineName
    }

    enableButton() {
        const $btn = this.$button
        $btn.disabled = false
        delete $btn.disabled
    }

    disableButton() {
        const $btn = this.$button
        $btn.disabled = true
    }
}