function slideshow() {
  return {
    albums: [],
    currentAlbum: '',
    photos: [],
    current: 0,
    playing: false,
    interval: 3,
    meta: { takenAt: '', model: '' },
    imageSrc: '',
    error: '',
    _timer: null,
    _abortCtrl: null,
    _preloadCache: new Map(),

    async init() {
      try {
        const res = await fetch('/api/albums')
        this.albums = await res.json() ?? []
      } catch {
        this.error = 'Failed to load albums.'
        return
      }
      if (this.albums.length > 0) {
        this.currentAlbum = this.albums[0].Key
        await this.loadPhotos()
      }
    },

    async loadPhotos() {
      clearInterval(this._timer)
      this.playing = false
      this.current = 0
      this.error = ''
      try {
        const res = await fetch('/api/albums/' + encodeURIComponent(this.currentAlbum))
        this.photos = await res.json() ?? []
      } catch {
        this.photos = []
        this.error = 'Failed to load photos.'
        return
      }
      this._revokePreloadCache()
      if (this.photos.length > 0) await this.loadImage()
    },

    photoUrl(token) {
      return '/photos/' + encodeURIComponent(this.currentAlbum) + '/' + encodeURIComponent(token)
    },

    _fetchPhoto(url, signal) {
      if (!this._preloadCache.has(url)) {
        const promise = fetch(url, { signal }).then(async res => {
          if (!res.ok) throw new Error('Network response was not ok');
          return {
            blob: await res.blob(),
            takenAt: res.headers.get('X-Photo-Taken-At'),
            model: res.headers.get('X-Photo-Model') || ''
          };
        }).catch(err => {
          this._preloadCache.delete(url);
          throw err;
        });
        this._preloadCache.set(url, promise);
      }
      return this._preloadCache.get(url);
    },

    async loadImage() {
      if (this.photos.length === 0) return

      if (this._abortCtrl) this._abortCtrl.abort()
      this._abortCtrl = new AbortController()
      const signal = this._abortCtrl.signal
      const target = this.current

      try {
        const url = this.photoUrl(this.photos[target])

        const data = await this._fetchPhoto(url, signal)

        if (signal.aborted) return

        this.error = ''
        this.meta = {
          takenAt: data.takenAt ? new Date(data.takenAt).toLocaleString() : '',
          model: data.model || '',
        }
        if (this.imageSrc) URL.revokeObjectURL(this.imageSrc)
        this.imageSrc = URL.createObjectURL(data.blob)
      } catch (e) {
        if (e.name === 'AbortError') return
        this.meta = { takenAt: '', model: '' }
        this.error = 'Failed to load image.'
      }
    },

    preloadAdjacent(index) {
      const targets = [
        (index + 1) % this.photos.length,
        (index - 1 + this.photos.length) % this.photos.length,
      ]
      for (const i of targets) {
        const url = this.photoUrl(this.photos[i])
        this._fetchPhoto(url).catch(() => {})
      }
    },

    _revokePreloadCache() {
      this._preloadCache.clear()
    },

    next() {
      this.current = (this.current + 1) % this.photos.length
      this.loadImage()
      this.preloadAdjacent(this.current)
      this.restartIfPlaying()
    },

    prev() {
      this.current = (this.current - 1 + this.photos.length) % this.photos.length
      this.loadImage()
      this.preloadAdjacent(this.current)
      this.restartIfPlaying()
    },

    goTo(i) {
      this.current = i
      this.loadImage()
      this.preloadAdjacent(i)
      this.restartIfPlaying()
    },

    togglePlay() {
      this.playing = !this.playing
      if (this.playing) {
        this._timer = setInterval(() => this.next(), this.interval * 1000)
      } else {
        clearInterval(this._timer)
      }
    },

    restartIfPlaying() {
      if (!this.playing) return
      clearInterval(this._timer)
      this._timer = setInterval(() => this.next(), this.interval * 1000)
    },
  }
}
