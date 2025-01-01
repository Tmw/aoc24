package main

import (
	"bufio"
	"container/list"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

type BlockType uint8

const (
	BlockTypeFree = BlockType(iota)
	BlockTypeFile
)

type Blocks struct {
	list *list.List
}

func (b *Blocks) String() string {
	var res strings.Builder
	for node := b.list.Front(); node != nil; node = node.Next() {
		b, ok := node.Value.(*Block)
		if !ok {
			continue
		}

		if b.Typ == BlockTypeFree {
			res.WriteString(strings.Repeat(".", int(b.Size)))
			continue
		}

		if b.Typ == BlockTypeFile {
			res.WriteString(strings.Repeat(fmt.Sprintf("%d", b.ID), int(b.Size)))
			continue
		}
	}

	return res.String()
}

func (b *Blocks) FindBlockOfTypeFromFront(typ BlockType) (*list.Element, *Block, int) {
	var idx int
	for node := b.list.Front(); node != nil; node = node.Next() {
		idx++
		n, ok := node.Value.(*Block)
		if !ok {
			fmt.Println("not correct typecast")
			return nil, nil, -1
		}

		if n.Typ != typ {
			continue
		}

		return node, n, idx
	}

	return nil, nil, -1
}

func (b *Blocks) FindBlockOfTypeFromBack(typ BlockType) (*list.Element, *Block, int) {
	var idx = b.list.Len()
	for node := b.list.Back(); node != nil; node = node.Prev() {
		idx--
		n, ok := node.Value.(*Block)
		if !ok {
			return nil, nil, -1
		}

		if n.Typ != typ {
			continue
		}

		return node, n, idx
	}

	return nil, nil, -1
}

func (b *Blocks) Checksum() int {
	var (
		idx int
		sum int
	)

	for node := b.list.Front(); node != nil; node = node.Next() {
		b, ok := node.Value.(*Block)
		if !ok {
			continue
		}

		if b.Typ != BlockTypeFile {
			continue
		}

		for range b.Size {
			sum += idx * b.ID
			idx++
		}
	}

	return sum
}

func (b *Blocks) Compact() {
	for {
		freeElm, freeBlock, freeIdx := b.FindBlockOfTypeFromFront(BlockTypeFree)
		fileElm, fileBlock, fileIdx := b.FindBlockOfTypeFromBack(BlockTypeFile)

		// exit the loop as soon as the first free block is past the first file block.
		if freeIdx > fileIdx {
			break
		}

		// We found the first available free block from the start,
		// and the first available file block from the end.
		// we'll need to handle the following scenario's:

		switch {
		case fileBlock.Size > freeBlock.Size:
			// File is bigger than free space
			// ------------------------------
			// - We insert a "file" block before the "free" block,
			// - We set the FileID to the same file ID we found near the end,
			// - We set the Size to the original file size - free space.
			// - We remove the original "free" block.
			b.list.InsertBefore(&Block{
				Typ:  BlockTypeFile,
				ID:   fileBlock.ID,
				Size: freeBlock.Size,
			}, freeElm)
			fileBlock.Size -= freeBlock.Size
			b.list.Remove(freeElm)

		case freeBlock.Size > fileBlock.Size:
			// Free space is bigger than file size
			// --------------------------------
			// - Move the "file" block before the "free" block
			// - Update the "free" block to its original size - file size.
			b.list.MoveBefore(fileElm, freeElm)
			freeBlock.Size -= fileBlock.Size

		case freeBlock.Size == fileBlock.Size:
			// File size and free space are exactly the same.
			// --------------------------------
			// - Move the "file" block before the "free" block
			// - Delete the free block.

			b.list.MoveBefore(fileElm, freeElm)
			b.list.Remove(freeElm)
		}
	}
}

type Block struct {
	ID   int
	Typ  BlockType
	Size uint8
}

func partOne(blocks Blocks) int {
	blocks.Compact()
	return blocks.Checksum()
}

func main() {
	blocks := parseInput(os.Stdin)

	start := time.Now()
	fmt.Println("answer part one =", partOne(blocks))
	fmt.Printf("part one took %+v\n", time.Since(start))
}

func parseInput(input io.Reader) Blocks {
	var (
		scanner = bufio.NewScanner(input)
		blocks  = Blocks{
			list: list.New(),
		}
		blockType = BlockTypeFile
		fileID    = 0
	)

	scanner.Split(bufio.ScanRunes)
	for scanner.Scan() {
		char := scanner.Text()
		if char == "\n" {
			continue
		}

		num, err := strconv.ParseUint(char, 10, 8)
		if err != nil {
			panic(fmt.Errorf("error converting %s to int: %w", char, err))
		}

		if blockType == BlockTypeFile {
			blocks.list.PushBack(&Block{ID: fileID, Typ: blockType, Size: uint8(num)})
			fileID++
			blockType = BlockTypeFree
			continue
		}

		if blockType == BlockTypeFree {
			blocks.list.PushBack(&Block{Typ: blockType, Size: uint8(num)})
			blockType = BlockTypeFile
			continue
		}
	}

	return blocks
}
