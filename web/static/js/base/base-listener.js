/**
 * Created by RTT.
 * Author: teocci@yandex.com on 2022-8월-22
 */
export default class BaseListener {
    constructor() {
        this._listeners = new Map()
    }

    listeners(event) {
        return this._listeners.get(event) ?? null
    }

    addListener(event, callback) {
        const callbacks = this.listeners(event) ?? []
        callbacks.push(callback)
        this._listeners.set(event, callbacks)
    }

    removeAllListeners() {
        this._listeners = new Map()
    }

    removeListener(event, callback = null) {
        const callbacks = this.listeners(event) ?? []
        if (isNil(callback)) {
            this._listeners.delete(event)
            return
        }

        // Remove specific handler
        for (const [i, cb] of callbacks.entries()) {
            if (cb === callback || cb.fn === callback) {
                callbacks.splice(i, 1)
                break
            }
        }

        // Remove event specific arrays for event types that no
        // one is subscribed for to avoid memory leak.
        if (callbacks.length === 0) {
            this._listeners.delete(event)
        }
    }

    once(event, callback) {
        const on = () => {
            this.removeListener(event, on)
            callback.apply(this, arguments)
        }
        on.fn = callback

        this.addListener(event, on)
    }

    callListener(event, e, ...params) {
        const callback = this.listeners(event) ?? null
        if (callback) callback(e, params)
    }

    emit(event) {
        const args = [...arguments].slice(1)
        const callbacks = this.listeners(event) ?? []

        for (const callback of callbacks) {
            callback.apply(this, args)
        }
    }

    hasListeners(event) {
        return !!this.listeners(event).length
    }
}