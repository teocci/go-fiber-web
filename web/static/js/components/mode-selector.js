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
    get buttons() {
        return this.buttonsMap.values() ?? []
    }

    /**
     * Returns the button's keys of the component.
     *
     * @returns {string[]}
     */
    get keys() {
        return this.buttonsMap.keys() ?? []
    }

    /**
     * Returns the mode buttons of the component.
     *
     * @returns {HTMLElement[]}
     */
    get modes() {
        return [...this.buttons].filter($btn => OPTION_MODE_KEYS.includes($btn.dataset.key))
    }

    /**
     * Returns the category buttons of the component.
     *
     * @returns {HTMLElement[]}
     */
    get categories() {
        return [...this.buttons].filter($btn => CATEGORY_KEYS.includes($btn.dataset.key))
    }

    initModeSelectorElements() {
        const $component = document.createElement('div')
        $component.classList.add(TAG, 'component-wrapper')

        const $bWrapper = document.createElement('div')
        $bWrapper.classList.add('buttons-list', 'bl-wrapper')

        for (const item of OPTION_MODES) {
            const $btn = this.createButtonElement(item)
            $component.appendChild($btn)
        }

        for (const item of CATEGORIES) {
            const $btn = this.createButtonElement(item)
            $component.appendChild($btn)
        }

        this.domWithHolderUpdate = $component
    }

    initModeSelectorListeners() {
        for (const $btn of this.buttons) {
            $btn.onclick = e => {
                this.onButtonClick(e, $btn)
            }
        }
    }

    createButtonElement(item) {
        const $btn = document.createElement('button')
        $btn.classList.add(`bl-${item.key}`, 'bl-item', 'hidden')
        $btn.dataset.key = item.key

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
                break
            case OPTION_MODE_KEY_CATEGORY:
                this.showCategories()
                break
            default:
                break
        }
    }
}