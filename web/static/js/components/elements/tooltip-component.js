/**
 * Created by Teocci.
 * Author: teocci@yandex.com on 2023-2월-17
 */
import BaseComponent from '../../base/base-component.js'

const TAG = 'tooltip'

export default class TooltipComponent extends BaseComponent {
    static TAG = TAG

    /** @typedef {Object} TooltipOptions
     * @property {string} content tooltip's content
     * @property {string} position tooltip's position
     */

    /** @type {?TooltipOptions} */
    options

    /** @type {HTMLDivElement} */
    $tooltip

    constructor(element, options) {
        super(element)

        this.options = cloner(options) ?? null
        if (isNil(this.options)) {
            this.options = {}
            this.options.content = this.dom.dataset.cooltip
            this.options.position = this.dom.dataset.cooltipPosition
        }
        if (isNil(this.options.content)) throw new Error('TooltipComponent: content is null')

        this.$tooltip = document.querySelector('.tooltip')
        if (isNil(this.$tooltip)) {
            this.$tooltip = document.createElement('div')
            this.$tooltip.classList.add('tooltip', 'tooltip-arrow')
            document.body.appendChild(this.$tooltip)
        }

        this.position = this.options.position
        this.content = this.options.content

        this.dom.onmouseenter = e => {
            this.showTooltip(e)
        }
        this.dom.onmouseleave = this.hideTooltip.bind(this)
    }

    set content(value) {
        if (isNil(this.$tooltip)) return
        if (isNil(value)) throw new Error('Content is null')

        this.options.content = value
        this.$tooltip.textContent = this.options.content
    }

    set position(value) {
        this.options.position = value ?? 'top'
        this.$tooltip.style.position = 'absolute'
    }

    showTooltip(e) {
        if (isNil(this.$tooltip)) return

        this.$tooltip.style.left = `${e.pageX}px`
        this.$tooltip.style.top = `${e.pageY}px`

        this.$tooltip.classList.add('show')
    }

    hideTooltip() {
        if (isNil(this.$tooltip)) return
        this.$tooltip.classList.remove('show')
    }
}