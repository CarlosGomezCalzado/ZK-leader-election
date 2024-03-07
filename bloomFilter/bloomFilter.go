package bloomFilter

import (
	"hash/fnv"
)

// BloomFilter es la estructura principal del filtro de Bloom
type BloomFilter struct {
	size    uint
	hashFn1  func(data []byte) uint
	hashFn2  func(data []byte) uint
	bitmap []bool
}

// NewBloomFilter crea un nuevo filtro de Bloom
func NewBloomFilter(size uint) *BloomFilter {
	return &BloomFilter{
		size:    size,
		hashFn1: hashFn1,
		hashFn2: hashFn2,
		bitmap:  make([]bool, size),
	}
}

func SetStateBloom(bitmap []bool) *BloomFilter {
	return &BloomFilter{
		size:    500,
		hashFn1: hashFn1,
		hashFn2: hashFn2,
		bitmap:  bitmap,
	}
}

func (bf *BloomFilter) GetBitmap()[]bool{
	return bf.bitmap
}

// Add agrega un elemento al filtro de Bloom
func (bf *BloomFilter) Add(data []byte) {
	index1 := bf.hashFn1(data) % bf.size
	index2 := bf.hashFn2(data) % bf.size

	bf.bitmap[index1] = true
	bf.bitmap[index2] = true
}

// Contains verifica si un elemento está presente en el filtro de Bloom
func (bf *BloomFilter) Contains(data []byte) bool {
	index1 := bf.hashFn1(data) % bf.size
	index2 := bf.hashFn2(data) % bf.size

	return bf.bitmap[index1] && bf.bitmap[index2]
}

// hashFn1 es una función de hash básica
func hashFn1(data []byte) uint {
	hash := fnv.New32a()
	hash.Write(data)
	return uint(hash.Sum32())
}

// hashFn2 es otra función de hash básica
func hashFn2(data []byte) uint {
	hash := fnv.New32()
	hash.Write(data)
	return uint(hash.Sum32())
}