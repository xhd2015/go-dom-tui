package dom

func (c *Node) GetEventHandler(eventType string) EventHandler {
	if c.Props == nil {
		return nil
	}

	h := getPropHandler(c.Props, "on"+eventType)
	if h != nil {
		return h
	}
	if eventType == "keydown" {
		h := getPropHandler(c.Props, "onKeyDown")
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

	fn, ok := handler.(func(event *DOMEvent) interface{})
	if ok && fn != nil {
		return fn
	}

	fnNoRes, ok := handler.(func(event *DOMEvent))
	if ok && fnNoRes != nil {
		return func(event *DOMEvent) interface{} {
			fnNoRes(event)
			return nil
		}
	}

	h, ok := handler.(EventHandler)
	if ok && h != nil {
		return h
	}

	return nil
}
