package dom

func (c *Node) GetEventHandler(eventType EventType) EventHandler {
	if c.Props == nil {
		return nil
	}

	h := getPropHandler(c.Props, "on"+string(eventType))
	if h != nil {
		return h
	}
	if eventType == "keydown" {
		h := getPropHandler(c.Props, "onKeyDown")
		if h != nil {
			return h
		}
	}
	if eventType == "resize" {
		h := getPropHandler(c.Props, "onWindowResize")
		if h != nil {
			return h
		}
	}
	return nil
}

func getPropHandler(props Props, key string) EventHandler {
	handler, ok := props.Get(key)
	if !ok {
		return nil
	}

	fn, ok := handler.(func(event *DOMEvent))
	if ok && fn != nil {
		return fn
	}

	fnNoRes, ok := handler.(func(event *DOMEvent))
	if ok && fnNoRes != nil {
		return func(event *DOMEvent) {
			fnNoRes(event)
		}
	}

	h, ok := handler.(EventHandler)
	if ok && h != nil {
		return h
	}

	return nil
}
