package database

import (
	"github.com/google/uuid"
	"github.com/brane-app/types-library"

	"sort"
	"strconv"
	"testing"
	"time"
)

var (
	content types.Content = types.NewContent(
		"https://gastrodon.io/file/foobar",
		uuid.New().String(),
		"png",
		[]string{"some", "tags"},
		true, false,
	)
	writableContent map[string]interface{} = content.Map()
)

func contentOK(test *testing.T, data map[string]interface{}, have types.Content) {
	if data["id"].(string) != have.ID {
		test.Errorf("types.Content ID mismatch! have: %s, want: %s", have.ID, data["id"])
	}

	var tags []string
	var ok bool
	if tags, ok = data["tags"].([]string); !ok {
		tags = make([]string, 0)
	}

	var length = len(have.Tags)
	if length != len(tags) {
		test.Errorf("Tags mismatch! have: %v, want: %v", have.Tags, tags)
	}

	sort.Strings(tags)
	sort.Strings(have.Tags)

	for length != 0 {
		length--
		if tags[length] != have.Tags[length] {
			test.Errorf("Tags mismatch at %d! have: %v, want: %v", length, have.Tags, tags)
		}
	}
}

func populate(limit int) {
	var modified map[string]interface{}

	var index int = 0
	var now int64 = time.Now().Unix()
	for index != limit {
		modified = mapCopy(writableContent)
		modified["id"] = uuid.New().String()
		modified["created"] = now + int64(100*index)
		WriteContent(modified)

		index++
	}

	return
}

func populateAuthor(author string, limit int) {
	var now int64 = time.Now().Unix()
	var modified map[string]interface{}
	var index int = 0
	for index != limit {
		modified = mapCopy(writableContent)
		modified["created"] = now
		modified["id"] = uuid.New().String()
		modified["author"] = uuid.New().String()
		WriteContent(modified)

		modified = mapCopy(writableContent)
		modified["created"] = now
		modified["id"] = uuid.New().String()
		modified["author"] = author
		WriteContent(modified)

		modified = mapCopy(writableContent)
		modified["created"] = now
		modified["id"] = uuid.New().String()
		modified["author"] = uuid.New().String()
		WriteContent(modified)

		index++
	}

	return
}

func Test_WriteContent(test *testing.T) {
	var mods []map[string]interface{} = []map[string]interface{}{
		map[string]interface{}{},
		map[string]interface{}{
			"tags": []string{},
		},
		map[string]interface{}{
			"id":   "0",
			"mime": "' or 1=1; DROP TABLE users",
		},
	}

	var err error
	var mod map[string]interface{}
	for _, mod = range mods {
		mod = mapMod(writableContent, mod)
		if err = WriteContent(mod); err != nil {
			test.Fatal(err)
		}
	}
}

func Test_WriteContent_err(test *testing.T) {
	var mods []map[string]interface{} = []map[string]interface{}{
		map[string]interface{}{
			"id": "96910fdf-916b-4664-a38c-be42d0f2c0ce foobar",
		},
		map[string]interface{}{
			"foo": "bar",
		},
	}

	var mod map[string]interface{}
	mod = mapCopy(writableContent)
	delete(mod, "file_url")

	var err error
	if err = WriteContent(mod); err == nil {
		test.Errorf("data %+v produced no error!", mod)
	}

	for _, mod = range mods {
		mod = mapMod(writableContent, mod)
		if err = WriteContent(mod); err == nil {
			test.Errorf("data %+v produced no error!", mod)
		}
	}
}

func Test_DeleteContent(test *testing.T) {
	var id string = uuid.New().String()
	var writable map[string]interface{} = mapMod(
		writableContent,
		map[string]interface{}{"id": id},
	)

	var err error
	if err = WriteContent(writable); err != nil {
		test.Fatal(err)
	}

	if err = DeleteContent(id); err != nil {
		test.Fatal(err)
	}

	var exists bool
	if _, exists, err = ReadSingleContent(id); err != nil {
		test.Fatal(err)
	}

	if exists {
		test.Errorf("deleted content %s still exists", id)
	}

	var tags []string
	if tags, err = getTags(id); err != nil {
		test.Fatal(err)
	}

	if len(tags) != 0 {
		test.Errorf("Tags %#v were not deleted for post %s", tags, id)
	}
}

func Test_ReadSingleContent(test *testing.T) {
	var modified map[string]interface{} = mapCopy(writableContent)
	modified["id"] = uuid.New().String()

	WriteContent(modified)

	var content types.Content = types.Content{}
	var exists bool
	var err error
	if content, exists, err = ReadSingleContent(modified["id"].(string)); err != nil {
		test.Fatal(err)
	}

	if !exists {
		test.Errorf("content of id %s does not exist!", modified["id"])
	}

	contentOK(test, modified, content)
}

func Test_ReadSingleContent_notags(test *testing.T) {
	var modified map[string]interface{} = mapCopy(writableContent)
	modified["id"] = uuid.New().String()
	modified["tags"] = []string{}

	WriteContent(modified)

	var content types.Content = types.Content{}
	var exists bool
	var err error
	if content, exists, err = ReadSingleContent(modified["id"].(string)); err != nil {
		test.Fatal(err)
	}

	if !exists {
		test.Errorf("content of id %s does not exist!", modified["id"])
	}

	contentOK(test, modified, content)

	if content.Tags == nil {
		test.Errorf("tags for %s are nil instead of empty!", content.ID)
	}
}

func Test_ReadSingleContent_ManyTags(test *testing.T) {
	var modified map[string]interface{} = mapCopy(writableContent)
	modified["id"] = uuid.New().String()

	var count int = 255
	var tags []string = make([]string, count)

	var index int = 0
	for index != count {
		tags[index] = "some_" + strconv.Itoa(index)
		index++
	}

	modified["tags"] = tags

	WriteContent(modified)

	var content types.Content = types.Content{}
	var exists bool
	var err error
	if content, exists, err = ReadSingleContent(modified["id"].(string)); err != nil {
		test.Fatal(err)
	}

	if !exists {
		test.Errorf("content of id %s does not exist!", modified["id"])
	}

	contentOK(test, modified, content)
}

func Test_ReadSingleContent_NotExists(test *testing.T) {
	var id string = uuid.New().String()

	var content types.Content
	var exists bool
	var err error
	if content, exists, err = ReadSingleContent(id); err != nil {
		test.Fatal(err)
	}

	if exists {
		test.Errorf("Query for nonexisting id got %+v", content)
	}
}

func Test_ReadManyContent(test *testing.T) {
	EmptyTable(CONTENT_TABLE)
	populate(30)

	var count int = 10
	var content []types.Content
	var size int
	var err error
	if content, size, err = ReadManyContent("", count); err != nil {
		test.Fatal(err)
	}

	if len(content) != count {
		test.Errorf("block is wrong size! have: %d, want: %d", len(content), count)
	}

	if len(content) != size {
		test.Errorf("block size mismatch! have: %d, want: %d", len(content), size)
	}
}

func Test_ReadManyContent_order(test *testing.T) {
	EmptyTable(CONTENT_TABLE)
	populate(30)

	var content []types.Content
	var err error
	if content, _, err = ReadManyContent("", 10); err != nil {
		test.Fatal(err)
	}

	var index int
	for index = range content[1:] {
		if content[index].Created < content[index+1].Created {
			test.Errorf("created out of order! this: %d, next: %d", content[index].Created, content[index+1].Created)
		}
	}
}

func Test_ReadManyContent_after(test *testing.T) {
	EmptyTable(CONTENT_TABLE)
	populate(30)

	var count, offset int = 10, 5
	var first, second []types.Content
	var err error
	if first, _, err = ReadManyContent("", count); err != nil {
		test.Fatal(err)
	}

	if second, _, err = ReadManyContent(first[offset].ID, count); err != nil {
		test.Fatal(err)
	}

	var index int
	var single types.Content
	for index, single = range first[offset+1:] {
		if second[index].ID != single.ID {
			test.Errorf("IDs not aligned! have: %s, want: %s", second[index].ID, single.ID)
		}
	}
}

func Test_ReadManyContent_afterNothing(test *testing.T) {
	EmptyTable(CONTENT_TABLE)
	populate(30)

	var content []types.Content
	var err error
	if content, _, err = ReadManyContent("foobar", 10); err != nil {
		test.Fatal(err)
	}

	if len(content) != 0 {
		test.Errorf("read after nonexisting id got %v", content)
	}
}

func Test_ReadManyContent_emptyTags(test *testing.T) {
	EmptyTable(CONTENT_TABLE)
	populate(1)

	var modified map[string]interface{} = mapCopy(writableContent)
	modified["id"] = "empty_tags"
	modified["tags"] = make([]string, 1)
	WriteContent(modified)

	modified = mapCopy(writableContent)
	modified["id"] = "empty_literal_tags"
	modified["tags"] = []string{}
	WriteContent(modified)

	var sizeExpect int = 3
	var content []types.Content
	var size int
	var err error
	if content, size, err = ReadManyContent("", sizeExpect); err != nil {
		test.Fatal(err)
	}

	if sizeExpect != size {
		test.Errorf("size mismatch! have: %d, want: %d", size, sizeExpect)
	}

	var single types.Content
	for _, single = range content {
		if single.Tags == nil {
			test.Errorf("%s has nil tags!", single.ID)
		}
	}
}

func Test_ReadManyContent_withTags(test *testing.T) {
	EmptyTable(CONTENT_TABLE)

	var modified map[string]interface{} = mapCopy(writableContent)
	modified["id"] = "with_tags"
	modified["tags"] = []string{"foo"}
	WriteContent(modified)

	var content []types.Content
	var err error
	if content, _, err = ReadManyContent("", 1); err != nil {
		test.Fatal(err)
	}

	if content[0].Tags[0] != "foo" {
		test.Errorf("tags mismatch! have: %v, want: %v", content[0].Tags, []string{"foo"})
	}
}

func Test_ReadManyContent_fewer(test *testing.T) {
	EmptyTable(CONTENT_TABLE)

	var population int = 5
	populate(population)

	var size int
	var err error
	if _, size, err = ReadManyContent("", 20); err != nil {
		test.Fatal(err)
	}

	if size != population {
		test.Errorf("size mismatch! have: %d, want: %d", size, population)
	}
}

func Test_ReadAuthorContent(test *testing.T) {
	EmptyTable(CONTENT_TABLE)

	var author string = uuid.New().String()
	populateAuthor(author, 20)

	var count int = 5
	var content []types.Content
	var size int
	var err error
	if content, size, err = ReadAuthorContent(author, "", count); err != nil {
		test.Fatal(err)
	}

	if len(content) != count {
		test.Errorf("block is wrong size! have: %d, want: %d", len(content), count)
	}

	if len(content) != size {
		test.Errorf("block size mismatch! have: %d, want: %d", len(content), size)
	}

	var single types.Content
	for _, single = range content {
		if single.Author != author {
			test.Errorf("author mismatch! have: %s, want: %s", single.Author, author)
		}
	}
}

func Test_ReadAuthorContent_after(test *testing.T) {
	EmptyTable(CONTENT_TABLE)

	var author string = uuid.New().String()
	populateAuthor(author, 20)

	var count, offset int = 10, 5
	var first, second []types.Content
	var err error
	if first, _, err = ReadAuthorContent(author, "", count); err != nil {
		test.Fatal(err)
	}

	if second, _, err = ReadAuthorContent(author, first[offset].ID, count); err != nil {
		test.Fatal(err)
	}

	var index int
	var single types.Content
	for index, single = range first[offset+1:] {
		if single.ID != second[index].ID {
			test.Errorf("IDs not aligned! have: %s, want: %s", second[index].ID, single.ID)
		}
	}
}

func Test_ReadAuthorContent_fewer(test *testing.T) {
	EmptyTable(CONTENT_TABLE)

	var author string = uuid.New().String()
	var population int = 5
	populateAuthor(author, population)

	var content []types.Content
	var size int
	var err error
	if content, size, err = ReadAuthorContent(author, "", 30); err != nil {
		test.Fatal(err)
	}

	if len(content) != population {
		test.Errorf("too high population! have: %d, want: %d", len(content), population)
	}

	if size != population {
		test.Errorf("too large size! have: %d, want: %d", size, population)
	}

	var single types.Content
	for _, single = range content {
		if single.Author != author {
			test.Errorf("author mismatch! have: %s, want: %s", single.Author, author)
		}
	}
}
