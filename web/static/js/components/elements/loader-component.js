/**
 * Created by Teocci.
 * Author: teocci@yandex.com on 2023-2월-17
 */
import BaseComponent from '../../base/base-component.js'

const TAG = 'loader'

export default class LoaderComponent extends BaseComponent {
    static TAG = TAG

    constructor() {
        super()

        this.dom = document.querySelector('.loader')
        if (isNil(this.dom)) {
            this.dom = document.createElement('div')
            this.dom.classList.add('loader')

            const $spinner = document.createElement('div')
            $spinner.classList.add('loader-spinner')

            this.dom.appendChild($spinner)

            document.body.appendChild(this.dom)
        }

        this.dom = this.dom
    }

    startLoader() {
        if (isNil(this.dom)) return

        this.dom.style.display = 'flex'
    }

    stopLoader() {
        if (isNil(this.dom)) return

        this.dom.style.display = 'none'
    }
}