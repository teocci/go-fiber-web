/**
 * Created by RTT.
 * Author: teocci@yandex.com on 2022-12월-22
 */
import BaseComponent from './base-component.js'
import ObservableObject from './observable-object.js'

const STATE_KEY_INIT = 'initialized'
const STATE_KEY_DATA_LOADED = 'data-loaded'
const STATE_KEY_DATA_EMPTY = 'data-empty'

export default class BaseTable extends BaseComponent {
    static STATE_KEY_INIT = STATE_KEY_INIT
    static STATE_KEY_DATA_LOADED = STATE_KEY_DATA_LOADED
    static STATE_KEY_DATA_EMPTY = STATE_KEY_DATA_EMPTY

    /** {ObservableObject} */
    _state

    /** {Function} */
    resolver

    /** {?Object} */
    data

    /** {HTMLElement} */
    $head

    /** {HTMLElement} */
    $table

    constructor(element) {
        super(element)

        this._state = new ObservableObject()
        this.state = STATE_KEY_INIT

        this.initBaseTableElements()
    }

    set state(state) {
        this._state.value = state
    }

    get state() {
        return this._state.value
    }

    set onStateChange(fn) {
        this._state.onchange = fn
    }

    initBaseTableElements() {
        const $wrapper = document.createElement('div')
        $wrapper.classList.add('module-wrapper')

        this.domWithHolderUpdate = $wrapper
    }
}