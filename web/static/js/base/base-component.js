/**
 * Created by RTT.
 * Author: teocci@yandex.com on 2022-8월-22
 */
import BaseListener from './base-listener.js'

export default class BaseComponent extends BaseListener {
    /** @type {HTMLElement} */
    $element

    /** @type {HTMLElement} */
    $placeholder

    constructor($element) {
        super()

        // TODO
        this.$element = $element ?? null
        this.$placeholder = $element ?? null
    }

    /**
     * Returns the holder element of the component.
     *
     * @return {HTMLElement}
     */
    get holder() {
        return this.$placeholder
    }

    /**
     * Returns the element of the component.
     *
     * @return {HTMLElement}
     */
    get dom() {
        return this.$element
    }

    get tag() {
        return this.constructor.TAG
    }

    set holder($holder) {
        this.$placeholder = $holder
    }

    set dom($element) {
        this.$element = $element
    }

    set domWithHolderUpdate($element) {
        this.dom = $element
        if (isNil(this.dom)) return
        if (isNil(this.holder)) return
        this.holder.append(this.dom)
    }

    set holderWithDomUpdate($element) {
        this.holder = $element
        if (isNil(this.holder)) return
        if (isNil(this.dom)) return
        this.holder.append(this.dom)
    }

    toggle(val) {
        const $element = this.dom
        $element.classList.toggle('hidden', val)
    }

    show() {
        const $element = this.dom
        $element.classList.remove('hidden')
    }

    hide() {
        const $element = this.dom
        $element.classList.add('hidden')
    }

    /**
     * Adds all arguments passed, except those already present.
     *
     * @param {...string} tokens
     */
    addClass(...tokens) {
        this.dom.classList.add(...tokens)
    }

    /**
     * Removes arguments passed, if they are present.
     *
     * @param {...string} tokens
     */
    removeClass(...tokens) {
        this.dom.classList.remove(...tokens)
    }

    /**
     * If force is not given, "toggles" tokenner, removing it if it's present
     * and adding it if it's not present.
     * If force is true, adds tokenner (same as add()).
     * If force is false, removes tokenner (same as remove()).
     *
     * Returns true if tokenner is now present, and false otherwise
     *
     * @param {string} token
     * @param {boolean} [force]
     * @return {boolean}
     */
    toggleClass(token, force) {
        this.dom.classList.toggle(token, force)
    }

    destroyChildren($element) {
        $element = $element ?? this.dom
        while ($element.firstChild) {
            const $lastChild = $element?.lastChild ?? null
            if (!isNil($lastChild)) $lastChild.remove()
        }
    }
}