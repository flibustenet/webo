// Copyright (c) 2023 William Dode. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for details.

package webo

import "fmt"

type ErrRedirect struct {
	URL string
}

func (e ErrRedirect) Error() string {
	return e.URL
}

func Redirect(url string, args ...any) ErrRedirect {
	return ErrRedirect{fmt.Sprintf(url, args...)}
}
