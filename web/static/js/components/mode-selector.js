/**
 * Created by Teocci.
 * Author: teocci@yandex.com on 2023-Aug-08
 */
import BaseComponent from '../base/base-component.js'

const TAG = 'mode-selector'

const OPTION_MODE_KEY_COMPANY = 'company'
const OPTION_MODE_KEY_CATEGORY = 'category'

const OPTION_MODE_BY_COMPANY = {
    key: OPTION_MODE_KEY_COMPANY,
    label: 'By Company',
    icon: 'fa-briefcase',
    url: '/seller/',
}

const OPTION_MODE_BY_CATEGORY = {
    key: OPTION_MODE_KEY_CATEGORY,
    label: 'By Category',
    icon: 'fa-tags',
}

const OPTION_MODE_KEYS = [
    OPTION_MODE_KEY_COMPANY,
    OPTION_MODE_KEY_CATEGORY,
]

const OPTION_MODES = [
    OPTION_MODE_BY_COMPANY,
    OPTION_MODE_BY_CATEGORY,
]

const CATEGORY_KEY_FEMALE = 'female'
const CATEGORY_KEY_MALE = 'male'
const CATEGORY_KEY_KIDS = 'kids'

const CATEGORY_FEMALE = {
    key: CATEGORY_KEY_FEMALE,
    uid: 9000,
    parent: 563,
    label: 'Женские ароматы',
    seo: 'Женская парфюмерия',
    url: '/catalog/krasota/parfyumeriya/zhenskie-aromaty',
    shard: 'beauty4',
    query: 'cat=9000',
    default: true,
}
const CATEGORY_MALE = {
    key: CATEGORY_KEY_MALE,
    uid: 9001,
    parent: 563,
    label: 'Мужские ароматы',
    seo: 'Мужская парфюмерия',
    url: '/catalog/krasota/parfyumeriya/muzhskie-aromaty',
    shard: 'beauty3',
    query: 'cat=9001',
}
const CATEGORY_KIDS = {
    key: CATEGORY_KEY_KIDS,
    uid: 9232,
    parent: 563,
    label: 'Детские ароматы',
    seo: 'Детская парфюмерия',
    url: '/catalog/krasota/parfyumeriya/detskie-aromaty',
    shard: 'beauty3',
    query: 'cat=9232',
}

const CATEGORY_KEYS = [
    CATEGORY_KEY_FEMALE,
    CATEGORY_KEY_MALE,
    CATEGORY_KEY_KIDS,
]

const CATEGORIES = [
    CATEGORY_FEMALE,
    CATEGORY_MALE,
    CATEGORY_KIDS,
]

export default class ModeSelector extends BaseComponent {
    static TAG = TAG

    /** @type {Map<string, HTMLElement>} */
    buttonsMap = new Map()

    /** @type {Map<string, HTMLElement>} */
    inputsMap = new Map()

    /** @type {HTMLElement} */
    $bWrapper

    /** @type {HTMLElement} */
    $rWrapper

    /** @type {HTMLElement} */
    $fWrapper

    constructor($element) {
        super($element)

        this.initModeSelectorElements()
        this.initModeSelectorListeners()

        console.log('ModeSelector', {buttonsMap: this.buttonsMap})

        this.showModes()
    }

    /**
     * Returns the buttons of the component.
     *
     * @returns {HTMLElement[]}
     */
    get modes() {
        return this.buttonsMap.values() ?? []
    }

    /**
     * Returns the button's keys of the component.
     *
     * @returns {string[]}
     */
    get modeKeys() {
        return this.buttonsMap.keys() ?? []
    }

    /**
     * Returns the category buttons of the component.
     *
     * @returns {HTMLElement[]}
     */
    get categories() {
        return this.inputsMap.values() ?? []
    }

    initModeSelectorElements() {
        const $component = document.createElement('div')
        $component.classList.add(TAG, 'component-wrapper')

        const $bWrapper = document.createElement('div')
        $bWrapper.classList.add('buttons-list', 'list-wrapper', 'cw-part')

        const $rWrapper = document.createElement('div')
        $rWrapper.classList.add('radios-list', 'list-wrapper', 'cw-part')

        const $fWrapper = document.createElement('div')
        $fWrapper.classList.add('filters-list', 'list-wrapper', 'cw-part')

        for (const item of OPTION_MODES) {
            const $btn = this.createButtonElement(item)
            $bWrapper.appendChild($btn)
        }

        for (const item of CATEGORIES) {
            const $btn = this.createRadioElement(item)
            $rWrapper.appendChild($btn)
        }

        $component.append($bWrapper, $rWrapper, $fWrapper)

        this.$bWrapper = $bWrapper
        this.$rWrapper = $rWrapper
        this.$fWrapper = $fWrapper

        this.domWithHolderUpdate = $component
    }

    initModeSelectorListeners() {
        for (const $btn of this.modes) {
            $btn.onclick = e => {
                this.onButtonClick(e, $btn)
            }
        }

        for (const $radio of this.categories) {
            const $input = $radio.querySelector('input')
            $input.onchange = e => {
                const uid = $input.value
                this.requestFilters('category', uid)
                this.onCategoryChange(e, $radio)
            }
        }
    }

    createButtonElement(item) {
        const $btn = document.createElement('button')
        $btn.classList.add(`bl-${item.key}`, 'bl-item', 'hidden')
        $btn.dataset.key = item.key
        $btn.dataset.uid = item.uid

        if (!isNil(item.icon)) {
            const $icon = document.createElement('div')
            $icon.classList.add('icon')

            const $i = document.createElement('i')
            $i.classList.add('fa-solid', item.icon)

            $btn.appendChild($i)
        }

        const $label = document.createElement('span')
        $label.classList.add('label')
        $label.textContent = item.label
        $btn.appendChild($label)

        this.buttonsMap.set(item.key, $btn)

        return $btn
    }

    createRadioElement(item) {
        const $radio = document.createElement('div')
        $radio.classList.add(`rl-${item.key}`, 'rl-item', 'hidden')
        $radio.dataset.key = item.key
        $radio.dataset.uid = item.uid

        const suid = `category-${item.uid}`

        const $input = document.createElement('input')
        $input.classList.add('cat-input', 'field-part')
        $input.type = 'radio'
        $input.name = 'category'
        $input.value = item.uid
        $radio.dataset.key = item.key
        $radio.dataset.uid = item.uid
        $input.id = suid
        if (item.default) $input.checked = true

        const $label = document.createElement('label')
        $label.classList.add('cat-label', 'field-part')
        $label.htmlFor = suid
        $label.textContent = item.label

        $radio.append($input, $label)

        this.inputsMap.set(item.key, $radio)

        return $radio
    }

    showModes() {
        for (const $btn of this.modes) {
            if (isNil($btn)) continue
            $btn.classList.remove('hidden')
        }
    }

    hideModes() {
        for (const $btn of this.modes) {
            if (isNil($btn)) continue
            $btn.classList.add('hidden')
        }
    }

    showCategories() {
        for (const $btn of this.categories) {
            if (isNil($btn)) continue
            $btn.classList.remove('hidden')
        }
    }

    hideCategories() {
        for (const $btn of this.categories) {
            if (isNil($btn)) continue
            $btn.classList.add('hidden')
        }
    }

    onButtonClick(e, $btn) {
        const key = $btn.dataset.key
        switch (key) {
            case OPTION_MODE_KEY_COMPANY:
                this.hideCategories()
                this.requestFilters('seller', pageInfo.sellerId)
                break
            case OPTION_MODE_KEY_CATEGORY:
                const uid = $btn.dataset.uid
                this.showCategories()
                this.requestFilters('category', uid)
                break
            default:
                break
        }
    }

    updateFilters(data) {
        if (isNil(data)) return

        const {filters} = data
        if (isNil(filters)) return

        const $select = document.createElement('select')
        $select.classList.add('selector')
        $select.id = 'filter-selector'

        for (const filter of filters) {
            const $option = this.createFilterElement(filter)
            $select.appendChild($option)
        }
    }

    createFilterElement(item) {
        const $option = document.createElement('option')
        $option.classList.add(`bl-${item.key}`, 'bl-item', 'hidden')
        $option.dataset.key = item.key
        $option.value = item.uid
        $option.textContent = item.label

        return $option
    }

    requestFilters(action, id) {
        const url = `/api/filters/${action}/${id}`

        console.log('requestFilters', {url})

        // this.fetchFilters(url).then(data => {
        //     if (isNil(data)) return
        //     this.updateFilters(data)
        // }).catch(e => {
        //     console.error(e)
        // })
    }

    async fetchFilters(url) {
        const response = await fetch(url)
        if (!response.ok) return null
        const data = await response.json()
        if (isNil(data)) return null

        return data
    }

    onCategoryChange(e, $radio) {
        const key = $radio.dataset.key
        const uid = $radio.dataset.uid
        console.log('onCategoryChange', {key, uid})
    }
}