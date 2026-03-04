function slideshow() {
  return {
    albums: [],
    currentAlbum: '',
    photos: [],
    current: 0,
    shuffle: false,
    playing: false,
    interval: 3,
    meta: { takenAt: '', model: '' },
    imageSrc: '',
    error: '',
    _timer: null,
    _abortCtrl: null,
    _preloadCache: new Map(),
    _touchStartX: 0,
    _touchStartY: 0,
    expanded: false,
    _exitHintTimer: null,
    showExitHint: false,

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
        const url = '/api/albums/' + encodeURIComponent(this.currentAlbum) + (this.shuffle ? '?shuffle=true' : '')
        const res = await fetch(url)
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

    handleTouchStart(e) {
      this._touchStartX = e.touches[0].clientX
      this._touchStartY = e.touches[0].clientY
    },

    handleTouchEnd(e) {
      if (!this.photos.length) return
      const dx = e.changedTouches[0].clientX - this._touchStartX
      const dy = e.changedTouches[0].clientY - this._touchStartY
      if (Math.abs(dx) < 50 || Math.abs(dy) > Math.abs(dx)) return
      if (dx < 0) this.next()
      else this.prev()
    },

    toggleExpand() {
      this.expanded = !this.expanded
      if (!this.expanded) this.showExitHint = false
    },

    handleImageWrapClick() {
      if (!this.expanded) return
      this.showExitHint = true
      clearTimeout(this._exitHintTimer)
      this._exitHintTimer = setTimeout(() => { this.showExitHint = false }, 1500)
    },

    handleKey(e) {
      if (!this.photos.length) return
      switch (e.key) {
        case 'ArrowRight': this.next(); break
        case 'ArrowLeft':  this.prev(); break
        case ' ':          e.preventDefault(); this.togglePlay(); break
        case 'Escape':     if (this.expanded) this.expanded = false; else if (this.playing) this.togglePlay(); break
        case 'f':          this.toggleExpand(); break
      }
    },
  }
}
