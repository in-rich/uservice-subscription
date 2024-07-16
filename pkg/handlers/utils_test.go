package handlers_test

import (
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

func RequireGRPCCodesEqual(t *testing.T, err error, code codes.Code) {
	if code != codes.OK {
		require.Error(t, err)
		st, ok := status.FromError(err)
		require.True(t, ok)
		require.Equal(t, code, st.Code())
	} else {
		require.NoError(t, err)
	}
}
