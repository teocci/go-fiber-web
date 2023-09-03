/**
 * Created by Teocci.
 * Author: teocci@yandex.com on 2023-Mar-01
 */
import BaseComponent from '../base/base-component.js'
import APIModule from '../modules/api-module.js'

const TAG = 'seller-details'

export default class SellerDetails extends BaseComponent {
    static TAG = TAG

    /** @typedef {Object} SellerInfo
     * @property {number} id
     * @property {string} name
     * @property {string} trademark
     * @property {string} legalAddress
     * @property {string} fineName
     * @property {string} ogrn
     * @property {boolean} isUnknown
     */

    /** @type {HTMLImageElement} */
    $logo

    /** @type {HTMLHeadingElement} */
    $name

    /** @type {HTMLDivElement} */
    $trademark

    /** @type {HTMLDivElement} */
    $legalAddress

    /** @type {HTMLButtonElement} */
    $button

    /** @type {?SellerInfo} */
    seller

    constructor(element) {
        super(element)

        this.seller = null

        this.initElements()
        this.initListeners()

        APIModule.requestSeller({sellerId: pageInfo.sellerId}, this.loadData.bind(this))
    }

    initElements() {
        const $component = document.createElement('div')
        $component.classList.add(TAG, 'component-wrapper')

        const $infoWrap = document.createElement('div')
        $infoWrap.classList.add('info-wrap', 'cw-item')

        const $logoWrap = document.createElement('div')
        $logoWrap.classList.add('logo-wrap')

        const $logo = document.createElement('div')
        $logo.classList.add('logo')

        const $img = document.createElement('img')
        $img.src = '/img/seller-empty-logo.jpg'
        $img.width = 120
        $img.height = 50

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

    loadData(seller) {
        this.seller = seller

        this.$name.textContent = seller.name
        this.$trademark.textContent = seller?.trademark
        this.$legalAddress.textContent = seller?.legalAddress

        this.$logo.src = `/seller/logo/${seller?.id}`
        this.$logo.alt = seller?.fineName
    }

    enableExportButton() {
        const $btn = this.$button
        $btn.disabled = false
        delete $btn.disabled
    }

    disableExportButton() {
        const $btn = this.$button
        $btn.disabled = true
    }

    hideExportButton() {
        const $btn = this.$button
        $btn.classList.add('hidden')
    }

    showExportButton() {
        const $btn = this.$button
        $btn.classList.remove('hidden')
    }
}