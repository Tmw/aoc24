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

func (b *Blocks) Clone() Blocks {
	res := Blocks{
		list: list.New(),
	}

	for node := b.list.Front(); node != nil; node = node.Next() {
		b, ok := node.Value.(*Block)
		if !ok {
			continue
		}

		res.list.PushBack(&Block{
			Typ:  b.Typ,
			ID:   b.ID,
			Size: b.Size,
		})
	}

	return res
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

func (b *Blocks) FindBlockOfTypeAfter(typ BlockType, elm *list.Element) (*list.Element, *Block, int) {
	var idx int
	for node := elm; node != nil; node = node.Next() {
		idx++
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

func (b *Blocks) FindBlockOfTypeBefore(typ BlockType, elm *list.Element) (*list.Element, *Block, int) {
	var idx = b.list.Len()
	for node := elm; node != nil; node = node.Prev() {
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
			// fmt.Println("typecast failed", node.Value)
			continue
		}

		if b.Typ == BlockTypeFree {
			// fmt.Println("block type not file, skipping")
			idx += int(b.Size)
			continue
		}

		for range b.Size {
			sum += idx * b.ID
			idx++
		}
	}

	return sum
}

func (b *Blocks) CompactWithoutFragmentation() {
	var (
		fileElm   = b.list.Back()
		fileBlock *Block
	)

	for {
		fileElm, fileBlock, _ = b.FindBlockOfTypeBefore(BlockTypeFile, fileElm)
		if fileElm == nil {
			break
		}

		var freeIdx int
		for node := b.list.Front(); node != nil; node = node.Next() {
			freeIdx++

			// break off search for a suitable free spot once we hit the
			// fileElm - this means we're in the middle. Since file blocks should only move left,
			// we should break off the search here.
			if node == fileElm {
				fileElm = fileElm.Prev()
				break
			}

			freeBlock, ok := node.Value.(*Block)
			if !ok {
				continue
			}

			if freeBlock.Typ != BlockTypeFree {
				continue
			}

			// file bigger than free space, continue search
			if fileBlock.Size > freeBlock.Size {
				continue
			}

			if fileBlock.Size == freeBlock.Size {
				newNodePos := fileElm.Next()
				b.list.MoveBefore(fileElm, node)
				b.list.MoveBefore(node, newNodePos)
				fileElm = newNodePos
				break
			}

			if fileBlock.Size < freeBlock.Size {
				newNodePos := b.list.InsertBefore(&Block{
					Typ:  BlockTypeFree,
					Size: fileBlock.Size,
				}, fileElm)
				b.list.MoveBefore(fileElm, node)
				freeBlock.Size -= fileBlock.Size

				fileElm = newNodePos
				break
			}
		}
	}
}

func (b *Blocks) CompactWithFragmentation() {
	for {
		freeElm, freeBlock, freeIdx := b.FindBlockOfTypeAfter(BlockTypeFree, b.list.Front())
		fileElm, fileBlock, fileIdx := b.FindBlockOfTypeBefore(BlockTypeFile, b.list.Back())

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
			// - We move the original "free" block to the end
			b.list.InsertBefore(&Block{
				Typ:  BlockTypeFile,
				ID:   fileBlock.ID,
				Size: freeBlock.Size,
			}, freeElm)
			fileBlock.Size -= freeBlock.Size
			b.list.MoveToBack(freeElm)

		case freeBlock.Size > fileBlock.Size:
			// Free space is bigger than file size
			// --------------------------------
			// - Move the "file" block before the "free" block
			// - Update the "free" block to its original size - file size.
			// - append new free space towards the end
			b.list.MoveBefore(fileElm, freeElm)
			freeBlock.Size -= fileBlock.Size
			b.list.PushBack(&Block{
				Typ:  BlockTypeFree,
				Size: fileBlock.Size,
			})

		case freeBlock.Size == fileBlock.Size:
			// File size and free space are exactly the same.
			// --------------------------------
			// - Move the "file" block before the "free" block
			// - Move the "free" block towards the end
			b.list.MoveBefore(fileElm, freeElm)
			b.list.MoveToBack(freeElm)
		}
	}
}

type Block struct {
	ID   int
	Typ  BlockType
	Size uint8
}

func partOne(blocks Blocks) int {
	blocks.CompactWithFragmentation()
	return blocks.Checksum()
}

func partTwo(blocks Blocks) int {
	blocks.CompactWithoutFragmentation()
	return blocks.Checksum()
}

func main() {
	blocks := parseInput(os.Stdin)
	blocks2 := blocks.Clone()

	start := time.Now()
	fmt.Println("answer part one =", partOne(blocks))
	fmt.Printf("part one took %+v\n", time.Since(start))

	start = time.Now()
	fmt.Println("answer part two =", partTwo(blocks2))
	fmt.Printf("part two took %+v\n", time.Since(start))
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
