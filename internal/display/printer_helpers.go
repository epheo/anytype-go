package display

import (
	"context"
	"fmt"
	"strings"

	"github.com/olekukonko/tablewriter"

	"github.com/epheo/anytype-go/pkg/anytype"
)

// prefetchTypeInformation pre-fetches type information for objects in a space
func prefetchTypeInformation(objects []anytype.Object, client *anytype.Client, ctx context.Context) {
	// We only need to pre-fetch once per space
	spacesSeen := make(map[string]bool)

	for _, obj := range objects {
		if obj.SpaceID != "" && !spacesSeen[obj.SpaceID] {
			// This will populate the type cache for this space with a single API call
			client.GetTypeName(ctx, obj.SpaceID, "dummy-key")
			spacesSeen[obj.SpaceID] = true
		}
	}
}

// formatDisplayName formats the display name with icon
func formatDisplayName(obj anytype.Object) string {
	name := obj.Name
	if name == "" {
		name = "<no name>"
	}

	var iconStr string
	if obj.Icon != nil {
		if obj.Icon.Emoji != "" {
			iconStr = obj.Icon.Emoji
		} else if obj.Icon.Name != "" {
			iconStr = obj.Icon.Name
		}
	}

	// Use GetPaddedIcon to ensure consistent spacing regardless of icon presence or type
	paddedIcon := GetPaddedIcon(iconStr, iconFixedWidth)
	displayName := fmt.Sprintf("%s%s", paddedIcon, name)

	// Truncate name if too long
	if len(displayName) > maxNameLength {
		displayName = displayName[:maxNameLength-3] + "..."
	}

	return displayName
}

// getTypeNameString retrieves friendly type name
func getTypeNameString(obj anytype.Object, client *anytype.Client, ctx context.Context) string {
	var typeNameStr string
	if obj.Type != nil {
		typeNameStr = obj.Type.Name
		if client != nil && obj.Type.Key != "" {
			typeNameStr = client.GetTypeName(ctx, obj.SpaceID, obj.Type.Key)
		}
	} else {
		typeNameStr = "Unknown"
	}
	return typeNameStr
}

// formatTagsString formats tags as a string
func formatTagsString(obj anytype.Object) string {
	// Format tags
	tags := "-"
	if len(obj.Tags) > 0 {
		tags = strings.Join(obj.Tags, ", ")
		if len(tags) > maxTagsLength {
			tags = tags[:maxTagsLength-3] + "..."
		}
	}
	return tags
}

// getLayoutString gets the layout as a string
func getLayoutString(obj anytype.Object) string {
	layout := obj.Layout
	if layout == "" {
		layout = "-"
	}
	return layout
}

// appendObjectToTable appends object data to the table
func appendObjectToTable(
	table *tablewriter.Table,
	obj anytype.Object,
	client *anytype.Client,
	ctx context.Context,
) {
	displayName := formatDisplayName(obj)
	typeNameStr := getTypeNameString(obj, client, ctx)
	layout := getLayoutString(obj)
	tags := formatTagsString(obj)

	table.Append([]string{displayName, typeNameStr, layout, tags})
}
