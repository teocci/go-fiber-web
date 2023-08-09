/**
 * Created by Teocci.
 * Author: teocci@yandex.com on 2023-Aug-09
 */
import BaseComponent from './base-component.js'

const TAG = 'button'

const DEFAULT_OPTIONS = {
    type: TAG,
    key: undefined,
    id: undefined,
    label: undefined,
    disabled: undefined,
    mode: undefined,
    state: undefined,
    action: undefined,
    classes: undefined,
    href: undefined,
}

const optionsParser = (base, extended) => {
    if (isNil(base)) throw new Error('BaseButton: base is null')
    if (isNil(extended)) throw new Error('BaseButton: extended is null')

    const merged = {...base}
    merged.key = isNil(extended?.key) ? base?.key ?? undefined : extended.key
    merged.id = isNil(extended?.id) ? base?.id ?? undefined : extended.id
    merged.label = isNil(extended?.label) ? base?.label ?? undefined : extended.label
    merged.disabled = isNil(extended?.disabled) ? base?.disabled ?? undefined : extended.disabled
    merged.mode = isNil(extended?.mode) ? base?.mode ?? undefined : extended.mode
    merged.state = isNil(extended?.state) ? base?.state ?? undefined : extended.state
    merged.action = isNil(extended?.action) ? base?.action ?? merged.key : extended.action
    merged.classes = mergeArraysWD(base?.classes, extended?.classes)
    merged.href = isNil(extended?.href) ? base?.href ?? undefined : extended.href

    return merged
}

export default class BaseButton extends BaseComponent {
    static TAG = TAG

    /** @typedef {Object} BaseButtonOptions
     * @property {string} type is always 'button'
     * @property {string} key to identify the button
     * @property {?string=} id HTML id
     * @property {string} label button's label
     * @property {?string=} disabled button is disabled if true
     * @property {?string=} mode button's mode
     * @property {?string=} state button's state
     * @property {?string=} action button's action
     * @property {?string[]=} classes button's classes
     * @property {?string=} href button's href used by the action
     */

    /** @type {BaseButtonOptions} */
    static DEFAULT_OPTIONS = DEFAULT_OPTIONS

    static optionsParser = optionsParser

    /** @type {BaseButtonOptions} */
    options = {}

    constructor($element, options) {
        super($element, options)

        this.parseBaseButtonOptions(options)

        this.initBCElement()
        this.loadBCClasses()
    }

    set onClick(f) {
        if (isNil(f)) throw new Error('InvalidParameter: f is null')
        if (!isFunction(f)) throw new Error('InvalidParameter: f is not a function')
        if (isNil(this.dom)) throw new Error('InvalidState: dom is null')

        this.dom.onclick = f
    }

    /**
     * Returns the default options for a base button
     *
     * @returns {BaseButtonOptions}
     */
    get defaultOptions() {
        const base = cloner(DEFAULT_OPTIONS)
        const extended = cloner(this.constructor.DEFAULT_OPTIONS)

        return BaseButton.optionsParser(base, extended)
    }

    get label() {
        if (isNil(this.dom)) return null
        return this.dom?.textContent ?? null
    }

    get action() {
        if (isNil(this.dom)) return null
        return this.dom.dataset?.action ?? null
    }

    get state() {
        if (isNil(this.dom)) return null
        return this.dom.dataset?.state ?? null
    }

    get mode() {
        if (isNil(this.dom)) return null
        return this.dom.dataset?.state ?? null
    }

    get disabled() {
        if (isNil(this.dom)) return false
        return this.dom?.disabled ?? false
    }

    set label(value) {
        if (isNil(this.dom)) return null
        this.dom.textContent = value
    }

    set action(value) {
        if (isNil(this.dom)) return null
        this.dom.dataset.action = value
    }

    set state(value) {
        if (isNil(this.dom)) return null
        this.dom.dataset.state = value
    }

    set mode(value) {
        if (isNil(this.dom)) return null
        this.dom.dataset.mode = value
    }

    set disabled(value) {
        if (isNil(this.dom)) return null
        this.dom.disabled = value
    }

    initBCElement() {
        this.domWithHolderUpdate = this.createButton()
    }

    parseBaseButtonOptions(options) {
        this.options = BaseButton.optionsParser(this.defaultOptions, options)
    }

    loadBCClasses() {
        if (isNil(this.options)) return

        const classes = this.options.classes
        if (isNil(classes)) return
        if (!isArray(classes)) return
        if (classes.length < 1) return

        this.addClass(...classes)
    }

    createButton() {
        const options = this.options

        const $btn = document.createElement('button')
        $btn.classList.add('btn')

        this.updateElement($btn, options)

        if (!isNil(options.href)) $btn.onclick = e => {
            window.location.href = options.href
        }

        return $btn
    }

    updateElement($btn, options) {
        if (isNil($btn)) return
        if (isNil(options)) return

        if (!isNil(options?.mode)) $btn.dataset.mode = options.mode
        if (!isNil(options?.state)) $btn.dataset.state = options.state
        if (!isNil(options?.action)) $btn.dataset.action = options.action

        if (!isNil(options?.id)) $btn.id = options.id
        if (!isNil(options?.disabled)) $btn.disabled = options.disabled

        if (!isNil(options?.label)) $btn.textContent = options.label
    }

    updateOptions(options) {
        if (isNil(this.dom)) throw new Error('InvalidState: dom is null')
        if (isNil(options)) return

        this.parseBaseButtonOptions(options)
        this.updateElement(options)
    }

    click() {
        if (isNil(this.dom)) return
        this.dom.click()
    }
}