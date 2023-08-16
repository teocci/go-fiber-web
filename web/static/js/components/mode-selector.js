/**
 * Created by Teocci.
 * Author: teocci@yandex.com on 2023-Aug-08
 */
import BaseComponent from '../base/base-component.js'

const TAG = 'mode-selector'

const OPTION_MODE_KEY_COMPANY = 'company'
const OPTION_MODE_KEY_CATEGORY = 'category'

const SELECTOR_SUID_CATEGORY = 'category-selector'
const SELECTOR_SUID_FILTER = 'filter-selector'

const DEFAULT_FILTER_KEY = 'xsubject'

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

const CATEGORY_KEY_LIST = [
    CATEGORY_KEY_FEMALE,
    CATEGORY_KEY_MALE,
    CATEGORY_KEY_KIDS,
]

const CATEGORY_LIST = [
    CATEGORY_FEMALE,
    CATEGORY_MALE,
    CATEGORY_KIDS,
]

const CATEGORY_MAP = {
    [CATEGORY_KEY_FEMALE]: CATEGORY_FEMALE,
    [CATEGORY_KEY_MALE]: CATEGORY_MALE,
    [CATEGORY_KEY_KIDS]: CATEGORY_KIDS,
}

export default class ModeSelector extends BaseComponent {
    static TAG = TAG

    /** @type {Map<string, HTMLElement>} */
    buttonsMap = new Map()

    /** @type {Map<string, HTMLElement>} */
    inputsMap = new Map()

    /** @type {HTMLElement} */
    $bWrapper

    /** @type {HTMLElement} */
    $cWrapper

    /** @type {HTMLElement} */
    $fWrapper

    /** @type {?string} */
    currentCategory = null

    /** @type {?string} */
    currentFilter = null

    constructor($element) {
        super($element)

        this.initModeSelectorElements()
        this.initModeSelectorListeners()

        console.log('ModeSelector', {buttonsMap: this.buttonsMap})

        this.showModes()
        this.onModeCompany()
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

    get selectedCategory() {
        const $selector = document.getElementById(SELECTOR_SUID_CATEGORY)
        if (isNil($selector)) return null

        return $selector.value
    }

    get selectedFilter() {
        const $selector = document.getElementById(SELECTOR_SUID_FILTER)
        if (isNil($selector)) return null

        return $selector.value
    }

    set selectedCategory(v) {
        const $selector = document.getElementById(SELECTOR_SUID_CATEGORY)
        if (isNil($selector)) return

        console.log('selectedCategory', {v})

        $selector.value = v
        const event = new Event('change')
        $selector.dispatchEvent(event)
    }

    set selectedFilterKey(v) {
        const $selector = document.getElementById(SELECTOR_SUID_FILTER)
        if (isNil($selector)) return

        $selector.value = v
    }

    initModeSelectorElements() {
        const $component = document.createElement('div')
        $component.classList.add(TAG, 'component-wrapper')

        const $bWrapper = document.createElement('div')
        $bWrapper.classList.add('buttons-list', 'list-wrapper', 'cw-part')

        const $cWrapper = document.createElement('div')
        $cWrapper.classList.add('categories-list', 'list-wrapper', 'cw-part')

        const $fWrapper = document.createElement('div')
        $fWrapper.classList.add('filters-list', 'list-wrapper', 'cw-part')

        for (const item of OPTION_MODES) {
            const $btn = this.createButtonElement(item)
            $bWrapper.append($btn)
        }

        this.createCategorySelector($cWrapper)

        $component.append($bWrapper, $cWrapper, $fWrapper)

        this.$bWrapper = $bWrapper
        this.$cWrapper = $cWrapper
        this.$fWrapper = $fWrapper

        this.domWithHolderUpdate = $component
    }

    initModeSelectorListeners() {
        for (const $btn of this.modes) {
            $btn.onclick = e => {
                this.onButtonClick(e, $btn)
            }
        }

        const $cSelector = document.getElementById(SELECTOR_SUID_CATEGORY)
        $cSelector.onchange = e => {
            this.onCategoryChange(e)
        }
    }

    checkByKey(key) {
        if (isNil(key)) return null

        const $input = this.inputsMap.get(key)
        if (isNil($input)) throw new Error('InvalidAttribute: input is null')

        const event = new Event('change')
        $input.dispatchEvent(event)
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

            $btn.append($i)
        }

        const $label = document.createElement('span')
        $label.classList.add('label')
        $label.textContent = item.label
        $btn.append($label)

        this.buttonsMap.set(item.key, $btn)

        return $btn
    }

    createCategorySelector($wrapper) {
        const $select = document.createElement('select')
        $select.classList.add(SELECTOR_SUID_CATEGORY, 'selector', 'hidden')
        $select.id = SELECTOR_SUID_CATEGORY

        for (const item of CATEGORY_LIST) {
            const $option = this.createCategoryOption(item)
            $select.append($option)
        }

        $wrapper.append($select)
    }

    createCategoryOption(item) {
        const $option = document.createElement('option')
        $option.classList.add(`cat-${item.key}`, 'cat-option')
        $option.value = item.key
        $option.textContent = item.label

        return $option
    }

    createRadioElement(item) {
        const $radio = document.createElement('div')
        $radio.classList.add(`rl-${item.key}`, 'rl-item', 'hidden')
        $radio.dataset.key = item.key
        $radio.dataset.uid = item.uid
        $radio.dataset.default = item.default

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
        const $selector = document.getElementById(SELECTOR_SUID_CATEGORY)
        $selector.classList.remove('hidden')
    }

    hideCategories() {
        const $selector = document.getElementById(SELECTOR_SUID_CATEGORY)
        $selector.classList.add('hidden')
    }

    onButtonClick(e, $btn) {
        const key = $btn.dataset.key
        switch (key) {
            case OPTION_MODE_KEY_COMPANY:
                this.onModeCompany()
                break
            case OPTION_MODE_KEY_CATEGORY:
                this.onModeCategory()
                break
            default:
                throw new Error('InvalidAttribute: key is null')
        }
    }

    onModeCompany() {
        this.hideCategories()
        this.requestFilters('seller', pageInfo.sellerId)
    }

    onModeCategory() {
        this.selectedCategory = CATEGORY_KEY_FEMALE
        this.showCategories()
    }

    onCategoryChange(e) {
        const key = e.target.value
        console.log('onCategoryChange', {key})
        if (isNil(key)) throw new Error('InvalidAttribute: key is null')

        const item = CATEGORY_MAP[key]
        if (isNil(item)) throw new Error('InvalidAttribute: item is null')

        if (key === this.currentCategory) return

        this.currentCategory = key
        this.requestFilters('category', item.uid)

        console.log('onCategoryChange', {key, uid: item.uid})
    }

    updateFilters(data) {
        if (isNil(data)) return

        let datum
        for (const item of data) {
            if (item.key === DEFAULT_FILTER_KEY) datum = item
        }
        if (isNil(datum) || isNil(datum.items)) return
        this.destroyChildren(this.$fWrapper)

        console.log('updateFilters', {datum})

        const $select = document.createElement('select')
        $select.classList.add(SELECTOR_SUID_FILTER, 'selector')
        $select.id = SELECTOR_SUID_FILTER

        const all = {
            id: 1,
            name: 'Все',
            count: datum.items.reduce((acc, cur) => acc + cur.count, 0),
        }
        datum.items = [all, ...datum.items]

        for (const filter of datum.items) {
            const $option = this.createFilterElement(filter)
            $select.append($option)
        }

        this.$fWrapper.append($select)
    }

    createFilterElement(item) {
        const $option = document.createElement('option')
        $option.classList.add(`bl-${item.id}`, 'bl-item')
        $option.dataset.key = `${DEFAULT_FILTER_KEY}-${item.id}`
        $option.value = item.id
        $option.textContent = `${item.name} (${numberFormatter(item.count)})`

        return $option
    }

    requestFilters(action, id) {
        const url = `/api/v1/filters/${action}/${id}`

        console.log('requestFilters', {url})

        this.fetchFilters(url).then(d => {
            if (isNil(d)) throw new Error('Filters not found')

            const {data} = d
            if (isNil(data)) throw new Error('Filters data not found')

            this.updateFilters(data)
        }).catch(e => {
            console.error(e)
        })
    }

    async fetchFilters(url) {
        const response = await fetch(url)
        if (!response.ok) return null
        const data = await response.json()
        if (isNil(data)) return null

        return data
    }
}