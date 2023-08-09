/**
 * Created by Teocci.
 * Author: teocci@yandex.com on 2023-2월-17
 */
import BaseButton from '../../base/base-button.js'

const TAG = 'button'

const DEFAULT_OPTIONS = {
    type: TAG,
    key: undefined,
    id: undefined,
    label: undefined,
    action: undefined,
    classes: undefined,
    href: undefined,
}

export default class ButtonComponent extends BaseButton {
    static TAG = TAG

    /** @type {BaseButtonOptions} */
    static DEFAULT_OPTIONS = DEFAULT_OPTIONS
}