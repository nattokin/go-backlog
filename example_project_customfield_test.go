package backlog_test

import (
	"context"
	"fmt"

	backlog "github.com/nattokin/go-backlog"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
)

var (
	// ProjectCustomFieldService
	doerProjectCustomFieldAll            = newMockDoer(fixture.CustomField.ListJSON)
	doerProjectCustomFieldCreate         = newMockDoer(fixture.CustomField.SingleJSON)
	doerProjectCustomFieldUpdate         = newMockDoer(fixture.CustomField.SingleJSON)
	doerProjectCustomFieldDelete         = newMockDoer(fixture.CustomField.SingleJSON)
	doerProjectCustomFieldAddListItem    = newMockDoer(fixture.CustomField.SingleJSON)
	doerProjectCustomFieldUpdateListItem = newMockDoer(fixture.CustomField.SingleJSON)
	doerProjectCustomFieldDeleteListItem = newMockDoer(fixture.CustomField.SingleJSON)
)

func ExampleProjectCustomFieldService_All() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerProjectCustomFieldAll),
	)

	fields, _ := c.Project.CustomField.All(context.Background(), "TEST")
	fmt.Printf("ID: %d, Name: %s\n", fields[0].ID, fields[0].Name)
	// Output:
	// ID: 1, Name: Sprint
}

func ExampleProjectCustomFieldService_Create() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerProjectCustomFieldCreate),
	)

	field, _ := c.Project.CustomField.Create(context.Background(), "TEST", backlog.CustomFieldTypeText, "Sprint")
	fmt.Printf("ID: %d, Name: %s\n", field.ID, field.Name)
	// Output:
	// ID: 1, Name: Sprint
}

func ExampleProjectCustomFieldService_Update() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerProjectCustomFieldUpdate),
	)

	field, _ := c.Project.CustomField.Update(
		context.Background(),
		"TEST",
		1,
		c.Project.CustomField.Option.WithName("Sprint Updated"),
	)
	fmt.Printf("ID: %d, Name: %s\n", field.ID, field.Name)
	// Output:
	// ID: 1, Name: Sprint
}

func ExampleProjectCustomFieldService_Delete() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerProjectCustomFieldDelete),
	)

	field, _ := c.Project.CustomField.Delete(context.Background(), "TEST", 1)
	fmt.Printf("ID: %d, Name: %s\n", field.ID, field.Name)
	// Output:
	// ID: 1, Name: Sprint
}

func ExampleProjectCustomFieldService_AddListItem() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerProjectCustomFieldAddListItem),
	)

	field, _ := c.Project.CustomField.AddListItem(context.Background(), "TEST", 1, "Item1")
	fmt.Printf("ID: %d, Name: %s\n", field.ID, field.Name)
	// Output:
	// ID: 1, Name: Sprint
}

func ExampleProjectCustomFieldService_UpdateListItem() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerProjectCustomFieldUpdateListItem),
	)

	field, _ := c.Project.CustomField.UpdateListItem(context.Background(), "TEST", 1, 10, "Item1 Updated")
	fmt.Printf("ID: %d, Name: %s\n", field.ID, field.Name)
	// Output:
	// ID: 1, Name: Sprint
}

func ExampleProjectCustomFieldService_DeleteListItem() {
	c, _ := backlog.NewClient(
		"https://example.backlog.com",
		"token",
		backlog.WithDoer(doerProjectCustomFieldDeleteListItem),
	)

	field, _ := c.Project.CustomField.DeleteListItem(context.Background(), "TEST", 1, 10)
	fmt.Printf("ID: %d, Name: %s\n", field.ID, field.Name)
	// Output:
	// ID: 1, Name: Sprint
}
