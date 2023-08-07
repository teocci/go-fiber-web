/**
 * Created by Teocci.
 * Author: teocci@yandex.com on 2023-Aug-08
 */
import BaseComponent from '../base/base-component.js'

const TAG = 'mode-selector'

const OPTION_MODE_KEY_BY_COMPANY = 'by_company'
const OPTION_MODE_KEY_BY_CATEGORY = 'by_category'

const OPTION_MODE_BY_COMPANY = {
    key: OPTION_MODE_KEY_BY_COMPANY,
    label: 'By Company',
    icon: 'fa fa-building',
}

const OPTION_MODE_BY_CATEGORY = {
    key: OPTION_MODE_KEY_BY_CATEGORY,
    label: 'By Category',
    icon: 'fa fa-list',
}

const OPTION_MODES = [
    OPTION_MODE_BY_COMPANY,
    OPTION_MODE_BY_CATEGORY
]

const CATEGORY_KEY_FEMALE = 'female'
const CATEGORY_KEY_MALE = 'male'
const CATEGORY_KEY_KIDS = 'kids'

const CATEGORY_FEMALE = {
    key: CATEGORY_KEY_FEMALE,
    id: 9000,
    parent: 563,
    name: 'Женские ароматы',
    seo: 'Женская парфюмерия',
    url: '/catalog/krasota/parfyumeriya/zhenskie-aromaty',
    shard: 'beauty4',
    query: 'cat=9000',
}
const CATEGORY_MALE = {
    key: CATEGORY_KEY_MALE,
    id: 9001,
    parent: 563,
    name: 'Мужские ароматы',
    seo: 'Мужская парфюмерия',
    url: '/catalog/krasota/parfyumeriya/muzhskie-aromaty',
    shard: 'beauty3',
    query: 'cat=9001',
}
const CATEGORY_KIDS = {
    key: CATEGORY_KEY_KIDS,
    id: 9232,
    parent: 563,
    name: 'Детские ароматы',
    seo: 'Детская парфюмерия',
    url: '/catalog/krasota/parfyumeriya/detskie-aromaty',
    shard: 'beauty3',
    query: 'cat=9232',
}

const CATEGORIES = [
    CATEGORY_KEY_FEMALE,
    CATEGORY_MALE,
    CATEGORY_KIDS,
]

export default class ModeSelector extends BaseComponent {
    static TAG = TAG

    constructor($element) {
        super($element)

        this.initModeSelectorElements()
        this.initModeSelectorListeners()
    }

    initModeSelectorElements() {

    }

    initModeSelectorListeners() {

    }
}